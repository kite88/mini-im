/**
 * 消息提示音播放工具
 *
 * 移动端（尤其 iOS Safari）的自动播放限制，是「手机端经常没声音」的根因：
 * 1. 非用户手势触发的 audio.play() 会被直接拒绝（NotAllowedError），
 *    而 WebSocket 收到消息时触发的播放不属于用户手势；
 * 2. 每次 new Audio() 得到的新实例同样处于「未解锁」状态，必须再等一次手势，
 *    所以「每条消息 new 一个 Audio」在手机上几乎必然无声；
 * 3. 短时间内创建多个 Audio 实例，iOS 上会互相抢占音频通道，导致部分消息无声；
 * 4. 未预加载时第一条消息要现下载，声音会明显延迟或来不及播完。
 *
 * 处理方式：
 * - 全页只用一个 Audio 实例并提前预加载；
 * - 首次用户手势（触摸 / 点击 / 按键）时静默播放一次完成「解锁」；
 * - 解锁后由 WebSocket 消息触发的播放即可正常发声；未解锁时先试播，
 *   被拒则挂在手势上，解锁后补播这条提示音。
 *
 * 注意：iOS 的硬件静音开关（拨杆）会静音所有 HTMLAudio，网页无法绕过，
 * 属于系统行为，不是 bug。
 */
function SoundUtil(src) {
    this.src = src || '/web/assets/music/wt01.mp3'
    this.audio = null
    this.unlocked = false
    this.pending = false
    this.gestureEvents = ['touchstart', 'mousedown', 'keydown']
    this._onGesture = null
}

/** 创建并预加载唯一的 Audio 实例；同时挂上解锁手势监听 */
SoundUtil.prototype.init = function () {
    if (this.audio) {
        return this.audio
    }
    var audio = new Audio(this.src)
    audio.preload = 'auto'
    audio.load()
    this.audio = audio
    this.bindGesture()
    return audio
}

SoundUtil.prototype.bindGesture = function () {
    if (this._onGesture) {
        return
    }
    var _this = this
    this._onGesture = function () {
        _this.unlock()
    }
    for (var i = 0; i < this.gestureEvents.length; i++) {
        document.addEventListener(this.gestureEvents[i], this._onGesture, {passive: true})
    }
}

SoundUtil.prototype.unbindGesture = function () {
    if (!this._onGesture) {
        return
    }
    for (var i = 0; i < this.gestureEvents.length; i++) {
        document.removeEventListener(this.gestureEvents[i], this._onGesture)
    }
    this._onGesture = null
}

/** 在用户手势中静默播放一次，解锁后续的程序化播放 */
SoundUtil.prototype.unlock = function () {
    var _this = this
    var audio = this.init()
    if (this.unlocked) {
        this.unbindGesture()
        return
    }
    // 用 volume=0 而不是 muted=true：muted 属于「静音自动播放」，
    // Safari 允许它自动播放，但不会因此解锁之后的有声播放
    var volume = audio.volume
    audio.volume = 0
    var ret = audio.play()
    audio.volume = volume
    var onOk = function () {
        audio.pause()
        try {
            audio.currentTime = 0
        } catch (e) {
        }
        _this.unlocked = true
        _this.unbindGesture()
        if (_this.pending) {
            _this.pending = false
            _this.play()
        }
    }
    if (ret && ret.then) {
        ret.then(onOk).catch(function () {
            // 仍然被拒绝（极少见）：保留手势监听，下次交互再解锁
        })
    } else {
        onOk()
    }
}

/** 播放提示音；重复调用会从头重播，不会叠加多个实例 */
SoundUtil.prototype.play = function () {
    var _this = this
    var audio = this.init()
    if (!this.unlocked) {
        // 未解锁：先直接试一次（安卓 / 桌面多数允许）
        this.bindGesture()
    }
    try {
        audio.currentTime = 0
    } catch (e) {
    }
    var ret = audio.play()
    if (ret && ret.catch) {
        ret.catch(function () {
            // 被自动播放策略拦截：等首次用户手势解锁后补播
            _this.unlocked = false
            _this.pending = true
            _this.bindGesture()
        })
    }
}

var soundUtil = new SoundUtil()
