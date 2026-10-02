// Command e2e 是 mini-im 的端到端回归脚本：只依赖被测服务本身暴露的 HTTP / WebSocket 接口，
// 覆盖「注册 → 登录 → 搜索 → 加好友 → 收发消息 → 离线消息 → 已读」全链路。
//
// 用法（先启动服务）：
//
//	go run ./tools/e2e                        # 跑全部场景
//	go run ./tools/e2e -scenario=core         # 只跑主链路
//	go run ./tools/e2e -scenario=edge         # 只跑边界 / 并发场景
//	go run ./tools/e2e -addr=http://127.0.0.1:2580
//
// 注意：脚本会真实注册账号并写入数据，且断言依赖「账号此前不存在、彼此尚不是好友」，
// 因此请在干净库上运行；重复运行请先清库：
//
//	psql -c "TRUNCATE im_message, im_friend, im_friend_request, im_read_state, im_user;"
//	redis-cli -n 1 flushdb
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var (
	httpBase string
	wsBase   string
	failed   int
	total    int
)

// 注册规则相关的后端提示原文（断言直接比对，改规则时同步这里）
const (
	usernameRuleMsg = "用户名需为 6-20 位字母、数字、下划线或减号组合"
	passwordRuleMsg = "密码至少 6 位"
)

// ---------------------------------------------------------------
// 断言与 HTTP 工具
// ---------------------------------------------------------------

func check(name string, ok bool, detail string) {
	total++
	if ok {
		fmt.Printf("PASS  %s\n", name)
		return
	}
	failed++
	if detail == "" {
		fmt.Printf("FAIL  %s\n", name)
		return
	}
	fmt.Printf("FAIL  %s | %s\n", name, detail)
}

// api 发起一次 JSON 请求，返回解析后的响应体（StatusCode 为 float64）。
func api(method, path, token string, body any) map[string]any {
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, httpBase+path, rdr)
	if err != nil {
		return map[string]any{"StatusCode": -999.0, "Message": err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return map[string]any{"StatusCode": -999.0, "Message": err.Error()}
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{"StatusCode": -998.0, "Message": string(raw)}
	}
	return out
}

func httpGet(path string) (int, string, string) {
	resp, err := http.Get(httpBase + path)
	if err != nil {
		return 0, "", err.Error()
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, resp.Header.Get("Content-Type"), string(body)
}

// ---------------------------------------------------------------
// WebSocket 客户端
// ---------------------------------------------------------------

type wsMsg struct {
	Event       string         `json:"event"`
	Data        map[string]any `json:"data"`
	Sender      string         `json:"sender"`
	Recipient   string         `json:"recipient"`
	Content     string         `json:"content"`
	Time        string         `json:"time"`
	ContentType string         `json:"content_type"`
	Subjoin     struct {
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
	} `json:"subjoin"`
}

type wsClient struct {
	conn *websocket.Conn
	ch   chan wsMsg
}

func dialWS(token string) (*wsClient, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(wsBase+"/ws?token="+token, nil)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("%v (status=%d)", err, resp.StatusCode)
		}
		return nil, err
	}

	w := &wsClient{conn: conn, ch: make(chan wsMsg, 256)}
	go func() {
		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				close(w.ch)
				return
			}
			var m wsMsg
			if json.Unmarshal(raw, &m) == nil {
				w.ch <- m
			}
		}
	}()
	return w, nil
}

func (w *wsClient) send(recipient, content, contentType string) error {
	return w.conn.WriteJSON(map[string]any{
		"recipient":    recipient,
		"content":      content,
		"content_type": contentType,
	})
}

func (w *wsClient) expect(pred func(wsMsg) bool, timeout time.Duration) (wsMsg, bool) {
	deadline := time.After(timeout)
	for {
		select {
		case m, ok := <-w.ch:
			if !ok {
				return wsMsg{}, false
			}
			if pred(m) {
				return m, true
			}
		case <-deadline:
			return wsMsg{}, false
		}
	}
}

func (w *wsClient) drain(timeout time.Duration) {
	deadline := time.After(timeout)
	for {
		select {
		case _, ok := <-w.ch:
			if !ok {
				return
			}
		case <-deadline:
			return
		}
	}
}

func (w *wsClient) close() { w.conn.Close() }

// ---------------------------------------------------------------
// 小工具
// ---------------------------------------------------------------

func register(username, password, nickname string) map[string]any {
	return api("POST", "/api/register", "", map[string]string{
		"username": username, "password": password, "nickname": nickname,
	})
}

// login 返回 (token, uid)，失败时返回空串。
func login(username, password string) (string, string) {
	r := api("POST", "/api/login", "", map[string]string{"username": username, "password": password})
	if num(r, "StatusCode") != 0 {
		return "", ""
	}
	d := r["Data"].(map[string]any)
	return str(d["tokens"].(map[string]any), "token"), str(d["user"].(map[string]any), "id")
}

func dataArr(r map[string]any) []any {
	v, _ := r["Data"].([]any)
	return v
}

func str(m map[string]any, key string) string {
	s, _ := m[key].(string)
	return s
}

func num(m map[string]any, key string) float64 {
	f, _ := m[key].(float64)
	return f
}

func item(v any) map[string]any {
	m, _ := v.(map[string]any)
	return m
}

// ---------------------------------------------------------------
// 场景一：主链路（注册 → 登录 → 加好友 → 在线 / 离线消息 → 已读 → 登录态）
// ---------------------------------------------------------------

func runCore() {
	fmt.Println("\n========== 场景一：主链路 ==========")

	register("carol_01", "123456", "卡罗尔")
	register("bob_01", "123456", "鲍勃")

	tokenCarol, idCarol := login("carol_01", "123456")
	tokenBob, idBob := login("bob_01", "123456")
	check("登录两个账号", tokenCarol != "" && tokenBob != "", "")
	if tokenCarol == "" || tokenBob == "" {
		return
	}
	fmt.Printf("      carol=%s bob=%s\n", idCarol, idBob)

	// 注册校验与重复注册
	check("重复注册被拒绝", register("carol_01", "123456", "卡罗尔")["Message"] == "用户名已存在", "")
	check("密码过短被拒绝", register("shortpw", "12345", "")["Message"] == passwordRuleMsg, "")
	check("密码错误登录失败",
		api("POST", "/api/login", "", map[string]string{"username": "carol_01", "password": "wrong"})["Message"] == "用户名或密码错误", "")

	// 鉴权
	check("无 token 返回 -10", num(api("GET", "/api/friends", "", nil), "StatusCode") == -10, "")
	check("坏 token 返回 -10", num(api("GET", "/api/friends", "bad.token", nil), "StatusCode") == -10, "")
	if _, err := dialWS("bad-token"); err == nil {
		check("WebSocket 拒绝坏 token", false, "预期 401")
	} else {
		check("WebSocket 拒绝坏 token", strings.Contains(err.Error(), "401"), err.Error())
	}

	// 搜索
	users := dataArr(api("GET", "/api/users?keyword=bob_01", tokenCarol, nil))
	check("按账号搜索到 bob", len(users) == 1 && str(item(users[0]), "username") == "bob_01", fmt.Sprint(users))
	if len(users) == 1 {
		check("搜索结果标注关系 none", str(item(users[0]), "relation") == "none", fmt.Sprint(users[0]))
	}
	byNick := dataArr(api("GET", "/api/users?keyword="+urlEscape("鲍勃"), tokenCarol, nil))
	check("按昵称搜索到 bob", len(byNick) == 1, fmt.Sprint(byNick))
	check("搜不到自己", len(dataArr(api("GET", "/api/users?keyword=carol_01", tokenCarol, nil))) == 0, "")
	check("关键词为空返回 []", len(dataArr(api("GET", "/api/users?keyword=", tokenCarol, nil))) == 0, "")

	// WebSocket 连接与欢迎语。
	// 连接数断言取「建连前基线 + 1」而非绝对值：服务可能已有其他客户端（例如浏览器页面）在线，
	// 绝对值断言会因外部连接而失败，这里只验证本连接被正确计入。
	baseConn := num(api("GET", "/healthz", "", nil), "connCount")

	bobWS, err := dialWS(tokenBob)
	check("bob 建立 WebSocket 连接", err == nil, fmt.Sprint(err))
	if err != nil {
		return
	}
	defer bobWS.close()
	_, ok := bobWS.expect(func(m wsMsg) bool { return m.Content == "socket服务连接成功" }, 3*time.Second)
	check("bob 收到连接成功提示", ok, "3s 内未收到")

	if h := api("GET", "/healthz", "", nil); num(h, "connCount") == baseConn+1 {
		check("healthz 连接数随新连接 +1", true, "")
	} else {
		check("healthz 连接数随新连接 +1", false, fmt.Sprintf("基线=%v 当前=%s", baseConn, fmt.Sprint(h)))
	}

	// 加好友：申请 → 事件推送 → 同意 → 事件推送
	check("不能加自己为好友",
		api("POST", "/api/friend-requests", tokenCarol, map[string]string{"to_id": idCarol, "message": "self"})["Message"] == "不能添加自己为好友", "")
	check("申请不存在的用户",
		api("POST", "/api/friend-requests", tokenCarol, map[string]string{"to_id": "1", "message": "hi"})["Message"] == "用户不存在", "")

	applyResp := api("POST", "/api/friend-requests", tokenCarol, map[string]string{"to_id": idBob, "message": "卡罗尔想加你"})
	check("carol 发起好友申请", num(applyResp, "StatusCode") == 0, fmt.Sprint(applyResp))

	evt, ok := bobWS.expect(func(m wsMsg) bool { return m.Event == "friend_request" }, 3*time.Second)
	check("bob 收到 friend_request 事件", ok, "3s 内未收到")
	if ok {
		check("事件带申请人昵称", str(evt.Data, "nickname") == "卡罗尔", fmt.Sprint(evt.Data))
		check("事件带验证消息", str(evt.Data, "message") == "卡罗尔想加你", fmt.Sprint(evt.Data))
	}

	pending := dataArr(api("GET", "/api/friend-requests", tokenBob, nil))
	check("bob 有待处理申请", len(pending) == 1, fmt.Sprint(pending))
	if len(pending) == 0 {
		return
	}
	reqID := str(item(pending[0]), "id")
	if u := dataArr(api("GET", "/api/users?keyword=bob_01", tokenCarol, nil)); len(u) == 1 {
		check("申请后 relation 变为 requested", str(item(u[0]), "relation") == "requested", fmt.Sprint(u[0]))
	}

	carolWS, err := dialWS(tokenCarol)
	check("carol 建立 WebSocket 连接", err == nil, fmt.Sprint(err))
	if err != nil {
		return
	}
	defer carolWS.close()
	carolWS.expect(func(m wsMsg) bool { return m.Content == "socket服务连接成功" }, 3*time.Second)

	check("bob 同意申请", num(api("POST", "/api/friend-requests/"+reqID+"/accept", tokenBob, nil), "StatusCode") == 0, "")
	evt2, ok := carolWS.expect(func(m wsMsg) bool { return m.Event == "friend_accepted" }, 3*time.Second)
	check("carol 收到 friend_accepted 事件", ok, "3s 内未收到")
	if ok {
		check("事件带同意者昵称", str(evt2.Data, "nickname") == "鲍勃", fmt.Sprint(evt2.Data))
	}
	check("重复同意被拒绝", api("POST", "/api/friend-requests/"+reqID+"/accept", tokenBob, nil)["Message"] == "该好友申请已处理", "")
	check("成为好友后再申请被拒绝",
		api("POST", "/api/friend-requests", tokenCarol, map[string]string{"to_id": idBob, "message": "x"})["Message"] == "你们已经是好友了", "")

	friendsBob := dataArr(api("GET", "/api/friends", tokenBob, nil))
	friendsCarol := dataArr(api("GET", "/api/friends", tokenCarol, nil))
	check("bob 好友列表含 carol", len(friendsBob) == 1 && str(item(friendsBob[0]), "id") == idCarol, fmt.Sprint(friendsBob))
	check("carol 好友列表含 bob", len(friendsCarol) == 1 && str(item(friendsCarol[0]), "id") == idBob, fmt.Sprint(friendsCarol))
	if u := dataArr(api("GET", "/api/users?keyword=bob_01", tokenCarol, nil)); len(u) == 1 {
		check("relation 变为 friend", str(item(u[0]), "relation") == "friend", fmt.Sprint(u[0]))
	}
	check("待处理申请已清空", len(dataArr(api("GET", "/api/friend-requests", tokenBob, nil))) == 0, "")

	// 在线消息收发
	carolWS.send(idBob, "你好 bob 👋", "text")
	msg, ok := bobWS.expect(func(m wsMsg) bool { return m.Content == "你好 bob 👋" }, 3*time.Second)
	check("bob 收到在线消息", ok, "3s 内未收到")
	if ok {
		check("消息 sender 为 carol", msg.Sender == idCarol, msg.Sender)
		check("消息 recipient 为 bob", msg.Recipient == idBob, msg.Recipient)
		check("消息补充发送者昵称", msg.Subjoin.Nickname == "卡罗尔", msg.Subjoin.Nickname)
		check("消息补充发送者头像", msg.Subjoin.Avatar != "", msg.Subjoin.Avatar)
		check("content_type=text", msg.ContentType == "text", msg.ContentType)
		check("time 由服务端填充", msg.Time != "", msg.Time)
	}

	// 伪造 sender 无效；未知 content_type 被归一化为 text
	carolWS.conn.WriteJSON(map[string]any{
		"sender": "999", "recipient": idBob, "content": "spoof", "content_type": "bogus",
		"subjoin": map[string]string{"nickname": "hacker", "avatar": "x"},
	})
	spoof, ok := bobWS.expect(func(m wsMsg) bool { return m.Content == "spoof" }, 3*time.Second)
	check("服务端覆盖伪造的 sender / subjoin，并把未知 content_type 归一化为 text",
		ok && spoof.Sender == idCarol && spoof.Subjoin.Nickname == "卡罗尔" && spoof.ContentType == "text",
		fmt.Sprintf("%+v", spoof))

	// 自己发给自己 / 非法 recipient 被丢弃
	carolWS.send(idCarol, "self-message", "text")
	carolWS.send("abc", "bad-recipient", "text")
	_, got := bobWS.expect(func(m wsMsg) bool {
		return m.Content == "self-message" || m.Content == "bad-recipient"
	}, time.Second)
	check("非法 recipient 被丢弃", !got, "bob 不该收到该消息")

	// 落库与聊天记录
	time.Sleep(1500 * time.Millisecond)
	chatsBob := dataArr(api("GET", "/api/chats/"+idCarol, tokenBob, nil))
	check("bob 聊天记录 2 条", len(chatsBob) == 2, fmt.Sprintf("got=%d", len(chatsBob)))
	if len(chatsBob) == 2 {
		check("最新一条内容为 spoof", chatsBob[0].(map[string]any)["content"] == "spoof", fmt.Sprint(chatsBob[0]))
		check("对端消息 type=2", num(item(chatsBob[0]), "type") == 2, fmt.Sprint(chatsBob[0]))
		check("对端消息带 user 资料", str(item(item(chatsBob[0])["user"]), "id") == idCarol, fmt.Sprint(chatsBob[0]))
		check("最早一条内容正确", chatsBob[1].(map[string]any)["content"] == "你好 bob 👋", fmt.Sprint(chatsBob[1]))
	}
	chatsCarol := dataArr(api("GET", "/api/chats/"+idBob, tokenCarol, nil))
	if len(chatsCarol) == 2 {
		check("自己发出的消息 type=1", num(item(chatsCarol[0]), "type") == 1, fmt.Sprint(chatsCarol[0]))
	} else {
		check("carol 聊天记录 2 条", false, fmt.Sprintf("got=%d", len(chatsCarol)))
	}
	check("非法 fid 返回参数错误", api("GET", "/api/chats/abc", tokenBob, nil)["Message"] == "参数错误", "")

	// 会话列表与未读数
	convs := dataArr(api("GET", "/api/messages", tokenBob, nil))
	check("bob 会话列表 1 条", len(convs) == 1, fmt.Sprint(convs))
	if len(convs) == 1 {
		c := item(convs[0])
		check("会话 uid 为 carol", str(c, "uid") == idCarol, fmt.Sprint(c))
		check("会话展示最新一条", c["message"] == "spoof", fmt.Sprint(c))
		check("未读数 >= 2", num(c, "msg_count") >= 2, fmt.Sprint(c))
		check("会话带对端资料", str(item(c["user"]), "nickname") == "卡罗尔", fmt.Sprint(c))
	}
	check("标记已读", num(api("POST", "/api/read/"+idCarol, tokenBob, nil), "StatusCode") == 0, "")
	if convs2 := dataArr(api("GET", "/api/messages", tokenBob, nil)); len(convs2) == 1 {
		check("未读数清零", num(item(convs2[0]), "msg_count") == 0, fmt.Sprint(convs2[0]))
	}

	// 离线消息
	bobWS.close()
	time.Sleep(500 * time.Millisecond)
	carolWS.send(idBob, "offline-1", "text")
	carolWS.send(idBob, "offline-2", "text")
	time.Sleep(500 * time.Millisecond)

	bobWS2, err := dialWS(tokenBob)
	check("bob 重连", err == nil, fmt.Sprint(err))
	if err == nil {
		defer bobWS2.close()
		bobWS2.expect(func(m wsMsg) bool { return m.Content == "socket服务连接成功" }, 3*time.Second)
		m1, ok1 := bobWS2.expect(func(m wsMsg) bool { return strings.HasPrefix(m.Content, "offline-") }, 3*time.Second)
		m2, ok2 := bobWS2.expect(func(m wsMsg) bool { return strings.HasPrefix(m.Content, "offline-") }, 3*time.Second)
		check("离线消息补发", ok1 && ok2, fmt.Sprintf("%q %q", m1.Content, m2.Content))
		check("离线消息保持原顺序", m1.Content == "offline-1" && m2.Content == "offline-2", fmt.Sprintf("%q %q", m1.Content, m2.Content))

		time.Sleep(1500 * time.Millisecond)
		check("离线消息同样落库", len(dataArr(api("GET", "/api/chats/"+idCarol, tokenBob, nil))) == 4,
			fmt.Sprint(len(dataArr(api("GET", "/api/chats/"+idCarol, tokenBob, nil)))))
	}

	// 单点登录 / 登出
	check("重复登录挤掉旧 token", func() bool {
		login("carol_01", "123456")
		return num(api("GET", "/api/friends", tokenCarol, nil), "StatusCode") == -10
	}(), "")

	newToken, _ := login("carol_01", "123456")
	check("登出", num(api("POST", "/api/logout", newToken, nil), "StatusCode") == 0, "")
	check("登出后 token 失效", num(api("GET", "/api/friends", newToken, nil), "StatusCode") == -10, "")

	// 头像接口
	code, ctype, body := httpGet("/avatar?name=carol&label=" + urlEscape("卡罗尔"))
	check("头像接口返回 SVG", code == 200 && strings.Contains(ctype, "svg") && strings.Contains(body, "<svg"),
		fmt.Sprintf("%d %s", code, ctype))
}

// ---------------------------------------------------------------
// 场景二：边界与并发
// ---------------------------------------------------------------

func runEdge() {
	fmt.Println("\n========== 场景二：边界与并发 ==========")

	register("nick0_01", "123456", "")
	{
		r := api("POST", "/api/login", "", map[string]string{"username": "nick0_01", "password": "123456"})
		nick := ""
		if num(r, "StatusCode") == 0 {
			nick = str(item(r["Data"].(map[string]any)["user"]), "nickname")
		}
		check("昵称留空时回退为账号名", nick == "nick0_01", "nickname="+nick)
	}
	// 用户名规则：6-20 位字母、数字、下划线或减号（过短 / 过长 / 含非法字符都要拒绝）
	check("用户名不符合规则被拒绝",
		register("d", "123456", "短名")["Message"] == usernameRuleMsg &&
			register(strings.Repeat("a", 21), "123456", "")["Message"] == usernameRuleMsg &&
			register("bad name", "123456", "")["Message"] == usernameRuleMsg, "")

	register("dave_01", "123456", "戴夫")
	register("erin_01", "123456", "艾琳")
	register("frank_01", "123456", "弗兰克")

	tokenDave, idDave := login("dave_01", "123456")
	tokenErin, idErin := login("erin_01", "123456")
	tokenFrank, idFrank := login("frank_01", "123456")
	check("三个账号登录成功", tokenDave != "" && tokenErin != "" && tokenFrank != "", "")
	if tokenDave == "" || tokenErin == "" || tokenFrank == "" {
		return
	}

	// 验证消息截断
	api("POST", "/api/friend-requests", tokenDave, map[string]string{"to_id": idErin, "message": strings.Repeat("加", 150)})
	list := dataArr(api("GET", "/api/friend-requests", tokenErin, nil))
	check("erin 收到 dave 的申请", len(list) == 1, fmt.Sprint(list))
	reqDaveToErin := ""
	if len(list) == 1 {
		reqDaveToErin = str(item(list[0]), "id")
		n := len([]rune(str(item(list[0]), "message")))
		check("验证消息截断到 100 字", n == 100, fmt.Sprintf("len=%d", n))
	}

	// 非法 / 越权处理申请
	check("非法申请 ID 返回参数错误", api("POST", "/api/friend-requests/abc/accept", tokenErin, nil)["Message"] == "参数错误", "")
	check("不存在的申请返回不存在", api("POST", "/api/friend-requests/123456789/accept", tokenErin, nil)["Message"] == "好友申请不存在", "")
	api("POST", "/api/friend-requests", tokenDave, map[string]string{"to_id": idFrank, "message": "hi frank"})
	frankPending := dataArr(api("GET", "/api/friend-requests", tokenFrank, nil))
	check("dave 的申请到达 frank", len(frankPending) == 1, fmt.Sprint(frankPending))
	if len(frankPending) == 1 {
		check("不能处理别人的申请",
			api("POST", "/api/friend-requests/"+str(item(frankPending[0]), "id")+"/accept", tokenErin, nil)["Message"] == "好友申请不存在", "")
	}

	// 拒绝 → 重新申请
	check("erin 拒绝申请", num(api("POST", "/api/friend-requests/"+reqDaveToErin+"/reject", tokenErin, nil), "StatusCode") == 0, "")
	check("拒绝后待处理清空", len(dataArr(api("GET", "/api/friend-requests", tokenErin, nil))) == 0, "")
	check("重复拒绝返回已处理", api("POST", "/api/friend-requests/"+reqDaveToErin+"/reject", tokenErin, nil)["Message"] == "该好友申请已处理", "")

	reApply := api("POST", "/api/friend-requests", tokenDave, map[string]string{"to_id": idErin, "message": "again"})
	check("被拒后可重新申请", num(reApply, "StatusCode") == 0 && reApply["Data"].(map[string]any)["became_friend"] == false, fmt.Sprint(reApply))
	rel := dataArr(api("GET", "/api/friend-requests", tokenErin, nil))
	check("申请重新变为待处理", len(rel) == 1, fmt.Sprint(rel))
	if len(rel) == 1 {
		check("重复申请刷新验证消息", str(item(rel[0]), "message") == "again", fmt.Sprint(rel[0]))
	}

	// 互相申请 = 直接成为好友
	mutual := api("POST", "/api/friend-requests", tokenErin, map[string]string{"to_id": idDave, "message": "me too"})
	check("互相申请直接成为好友", num(mutual, "StatusCode") == 0 && mutual["Data"].(map[string]any)["became_friend"] == true, fmt.Sprint(mutual))
	dl := dataArr(api("GET", "/api/friends", tokenDave, nil))
	check("dave 好友列表含 erin", len(dl) == 1 && str(item(dl[0]), "id") == idErin, fmt.Sprint(dl))
	check("互相确认后待处理清空", len(dataArr(api("GET", "/api/friend-requests", tokenErin, nil))) == 0, "")

	// 图片消息 / 并发消息
	daveWS, err := dialWS(tokenDave)
	check("dave 建立 WebSocket 连接", err == nil, fmt.Sprint(err))
	if err != nil {
		return
	}
	defer daveWS.close()
	erinWS, err := dialWS(tokenErin)
	check("erin 建立 WebSocket 连接", err == nil, fmt.Sprint(err))
	if err != nil {
		return
	}
	defer erinWS.close()
	daveWS.drain(500 * time.Millisecond)
	erinWS.drain(500 * time.Millisecond)

	erinWS.send(idDave, `<img src="/assets/image/emoji/emoji_u1f382.png">`, "image")
	m, ok := daveWS.expect(func(m wsMsg) bool { return strings.HasPrefix(m.Content, "<img") }, 3*time.Second)
	check("图片消息（content_type=image）送达", ok && m.ContentType == "image", fmt.Sprintf("%+v", m))

	const bulk = 30
	for i := 0; i < bulk; i++ {
		erinWS.send(idDave, fmt.Sprintf("bulk-%02d", i), "text")
	}
	got := map[string]bool{}
	deadline := time.After(5 * time.Second)
	for len(got) < bulk {
		select {
		case msg, open := <-daveWS.ch:
			if !open {
				goto collected
			}
			if strings.HasPrefix(msg.Content, "bulk-") {
				got[msg.Content] = true
			}
		case <-deadline:
			goto collected
		}
	}
collected:
	check("30 条并发消息全部送达", len(got) == bulk, fmt.Sprintf("got=%d", len(got)))

	time.Sleep(2 * time.Second)
	persisted := len(dataArr(api("GET", "/api/chats/"+idErin, tokenDave, nil)))
	check("并发消息全部落库", persisted == bulk+1, fmt.Sprintf("persisted=%d want=%d", persisted, bulk+1))

	// 超长消息触发 8KB 读限制，只断开该连接
	bigWS, err := dialWS(tokenErin)
	check("erin 第二条连接建立", err == nil, fmt.Sprint(err))
	if err == nil {
		_ = bigWS.conn.WriteMessage(websocket.TextMessage,
			[]byte(`{"recipient":"`+idDave+`","content":"`+strings.Repeat("x", 9000)+`","content_type":"text"}`))
		closed := false
		for i := 0; i < 30; i++ {
			if _, _, err := bigWS.conn.ReadMessage(); err != nil {
				closed = true
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		check("超过 8KB 的消息断开连接", closed, "连接仍存活")
		bigWS.close()
	}
	if h := api("GET", "/healthz", "", nil); str(h, "status") == "ok" {
		check("服务仍健康", true, "")
	} else {
		check("服务仍健康", false, fmt.Sprint(h))
	}
	erinWS.send(idDave, "after-oversize", "text")
	_, ok = daveWS.expect(func(m wsMsg) bool { return m.Content == "after-oversize" }, 3*time.Second)
	check("其他连接不受影响", ok, "未收到消息")

	// token 续签
	rr := api("PUT", "/api/token", tokenFrank, nil)
	check("续签 token", num(rr, "StatusCode") == 0, fmt.Sprint(rr))
	newToken := ""
	if num(rr, "StatusCode") == 0 {
		newToken = str(rr["Data"].(map[string]any), "token")
	}
	check("续签后旧 token 失效", num(api("GET", "/api/friends", tokenFrank, nil), "StatusCode") == -10, "")
	check("续签后新 token 可用", num(api("GET", "/api/friends", newToken, nil), "StatusCode") == 0, "")
}

// urlEscape 只处理关键字里的少量中文，避免额外依赖
func urlEscape(s string) string {
	const hex = "0123456789ABCDEF"
	var b strings.Builder
	for _, c := range []byte(s) {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == '~' {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(hex[c>>4])
		b.WriteByte(hex[c&0x0f])
	}
	return b.String()
}

func main() {
	addr := flag.String("addr", "http://127.0.0.1:2580", "被测服务地址")
	scenario := flag.String("scenario", "all", "运行场景：all | core | edge")
	flag.Parse()

	httpBase = strings.TrimRight(*addr, "/")
	wsBase = strings.Replace(strings.Replace(httpBase, "https://", "wss://", 1), "http://", "ws://", 1)

	switch *scenario {
	case "core":
		runCore()
	case "edge":
		runEdge()
	case "all":
		runCore()
		runEdge()
	default:
		fmt.Printf("未知场景 %q\n", *scenario)
		os.Exit(2)
	}

	fmt.Printf("\n%d/%d 项通过", total-failed, total)
	if failed > 0 {
		fmt.Printf("，%d 项失败\n", failed)
		os.Exit(1)
	}
	fmt.Println("，全部通过 ✅")
}
