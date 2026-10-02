/**
 * 国际化（i18n）工具，支持语言：简体中文（zh-CN）、英文（en-US）
 *
 * 语言优先级：localStorage -> 浏览器语言 -> 默认 zh-CN
 *
 * 用法：
 *   var i18nUtil = new I18nUtil()
 *   i18nUtil.get()                                           // 读取当前语言
 *   i18nUtil.t('login.submit')                               // 取当前语言文案
 *   i18nUtil.t('index.account', {username: 'bln'}, 'en-US')  // 指定语言 + 占位符
 *   i18nUtil.set('en-US')                                    // 保存语言（同时更新 <html lang>）
 *
 * 页面上把语言放在 Vue data 上（lang）即可获得响应式切换，调用时显式传入：
 *   t(key, params) { return i18nUtil.t(key, params, this.lang) }
 */
function I18nUtil(lang) {
    this.lang = I18nUtil.normalize(lang || this.get())
}

I18nUtil.STORAGE_KEY = 'mini-im-lang';
I18nUtil.DEFAULT_LANG = 'zh-CN';
I18nUtil.SUPPORTED_LANGS = ['zh-CN', 'en-US'];

// 与 document.documentElement.lang 对应：zh-CN / en
I18nUtil.HTML_LANGS = {'zh-CN': 'zh-CN', 'en-US': 'en'};

I18nUtil.MESSAGES = {
    'zh-CN': {
        'page.loginTitle': 'mini-im登录',
        'page.indexTitle': 'mini-im',

        'common.tip': '提示',
        'common.confirm': '确认',
        'common.cancel': '取消',

        'settings.title': '设置',
        'settings.themeMode': '主题模式',
        'settings.fontSize': '字号',
        'settings.fontWeight': '字重',
        'settings.language': '多语言',

        'lang.title': '语言',
        'lang.zhTitle': '中文',
        'lang.enTitle': 'English',

        'theme.title': '主题',
        'theme.light': '浅色',
        'theme.system': '跟随系统',
        'theme.dark': '深色',

        'fontWeight.normal': '常规',
        'fontWeight.bold': '加粗',

        'login.username': '用户',
        'login.usernamePlaceholder': '请输入用户名',
        'login.password': '密码',
        'login.passwordPlaceholder': '请输入用户密码',
        'login.rememberMe': '请记住我',
        'login.clearPassword': '清除密码缓存',
        'login.submit': '登录',
        'login.usernameRequired': '请输入用户名',
        'login.passwordRequired': '请输入密码',
        'login.noAccount': '还没有账号？',

        'register.title': '注册账号',
        'register.toRegister': '立即注册',
        'register.toLogin': '返回登录',
        'register.username': '用户',
        'register.usernamePlaceholder': '6-20 位字母、数字、下划线或减号',
        'register.nickname': '昵称',
        'register.nicknamePlaceholder': '选填，默认与用户名相同',
        'register.password': '密码',
        'register.passwordPlaceholder': '请输入密码（至少 6 位）',
        'register.password2': '确认密码',
        'register.password2Placeholder': '请再次输入密码',
        'register.submit': '注册',
        'register.usernameRequired': '请输入用户名',
        'register.usernameRule': '用户名需为 6-20 位字母、数字、下划线或减号组合',
        'register.passwordRequired': '请输入密码',
        'register.passwordTooShort': '密码至少 6 位',
        'register.passwordNotMatch': '两次输入的密码不一致',
        'register.success': '注册成功，正在自动登录…',

        'index.account': '账号：{username}',
        'index.logout': '退出登录',
        'index.logoutConfirm': '确认退出登录吗？',
        'index.logoutTitle': '退出提示',
        'index.navMessage': '消息',
        'index.navFriend': '好友',
        'index.searchPlaceholder': '搜索',
        'index.searchSession': '搜索会话',
        'index.searchFriend': '搜索好友',
        'index.emptySession': '暂无会话',
        'index.emptyFriend': '暂无好友',
        'index.navNewFriend': '新的朋友',
        'index.addFriend': '添加好友',
        'index.search': '搜索',
        'index.searchUserPlaceholder': '输入账号或昵称',
        'index.greetingPlaceholder': '验证消息（选填）',
        'index.searching': '搜索中…',
        'index.keywordRequired': '请输入要搜索的账号或昵称',
        'index.emptySearchResult': '没有找到相关用户',
        'index.emptyFriendRequest': '暂无新的朋友',
        'index.add': '添加',
        'index.accept': '同意',
        'index.reject': '拒绝',
        'index.relationFriend': '已是好友',
        'index.relationRequested': '申请中',
        'index.relationIncoming': '待你处理',
        'index.relationBlocked': '已拉黑',
        'index.navBlacklist': '黑名单',
        'index.searchBlacklist': '搜索黑名单',
        'index.emptyBlacklist': '黑名单为空',
        'index.blacklistAdd': '拉黑',
        'index.blacklistRemove': '移出黑名单',
        'index.blacklistConfirm': '确定将 {nickname} 拉入黑名单吗？拉黑后对方发来的消息将被拦截。',
        'index.blacklistRemoveConfirm': '确定将 {nickname} 移出黑名单吗？',
        'index.messageBlocked': '{nickname} 已将你拉入黑名单，消息未能送达',
        'index.friendDelete': '删除',
        'index.friendDeleteConfirm': '确定删除好友 {nickname} 吗？你这边的好友列表、会话与聊天记录都会被清空，对方的记录不受影响；再次添加将是全新会话。',
        'index.messageNotFriend': '你们还不是好友，消息未能送达',
        'index.messageUnfriended': '对方已将你从好友中删除，消息未能送达',
        'index.defaultGreeting': '请求添加你为好友',
        'index.newFriendRequest': '{nickname} 请求添加你为好友',
        'index.friendAccepted': '{nickname} 已同意你的好友申请',
        'index.emojiMessage': '[动画表情]',
        'index.emoji': '表情',
        'index.send': '发送',
        'index.sendTitle': 'ctrl+Enter 发送',
        'index.selectTarget': '请选择聊天对象',
        'index.emptyContent': '请输入聊天内容',
        'index.warmTip': '温馨提示',
        'index.socketUnsupported': '您的浏览器不支持 socket'
    },
    'en-US': {
        'page.loginTitle': 'mini-im Sign In',
        'page.indexTitle': 'mini-im',

        'common.tip': 'Notice',
        'common.confirm': 'OK',
        'common.cancel': 'Cancel',

        'settings.title': 'Settings',
        'settings.themeMode': 'Theme mode',
        'settings.fontSize': 'Text size',
        'settings.fontWeight': 'Text weight',
        'settings.language': 'Language',

        'lang.title': 'Language',
        'lang.zhTitle': '中文',
        'lang.enTitle': 'English',

        'theme.title': 'Theme',
        'theme.light': 'Light',
        'theme.system': 'System',
        'theme.dark': 'Dark',

        'fontWeight.normal': 'Regular',
        'fontWeight.bold': 'Bold',

        'login.username': 'Username',
        'login.usernamePlaceholder': 'Please enter your username',
        'login.password': 'Password',
        'login.passwordPlaceholder': 'Please enter your password',
        'login.rememberMe': 'Remember me',
        'login.clearPassword': 'Clear saved password',
        'login.submit': 'Sign in',
        'login.usernameRequired': 'Please enter your username',
        'login.passwordRequired': 'Please enter your password',
        'login.noAccount': "Don't have an account?",

        'register.title': 'Create account',
        'register.toRegister': 'Sign up now',
        'register.toLogin': 'Back to sign in',
        'register.username': 'Username',
        'register.usernamePlaceholder': '6-20 letters, digits, underscore or hyphen',
        'register.nickname': 'Nickname',
        'register.nicknamePlaceholder': 'Optional, same as username by default',
        'register.password': 'Password',
        'register.passwordPlaceholder': 'Please enter your password (min 6 chars)',
        'register.password2': 'Confirm',
        'register.password2Placeholder': 'Please enter your password again',
        'register.submit': 'Sign up',
        'register.usernameRequired': 'Please enter your username',
        'register.usernameRule': 'Username must be 6-20 letters, digits, underscore or hyphen',
        'register.passwordRequired': 'Please enter your password',
        'register.passwordTooShort': 'Password must be at least 6 characters',
        'register.passwordNotMatch': 'The two passwords do not match',
        'register.success': 'Signed up successfully, signing in…',

        'index.account': 'Account: {username}',
        'index.logout': 'Sign out',
        'index.logoutConfirm': 'Are you sure you want to sign out?',
        'index.logoutTitle': 'Sign out',
        'index.navMessage': 'Messages',
        'index.navFriend': 'Contacts',
        'index.searchPlaceholder': 'Search',
        'index.searchSession': 'Search conversations',
        'index.searchFriend': 'Search contacts',
        'index.emptySession': 'No conversations yet',
        'index.emptyFriend': 'No contacts yet',
        'index.navNewFriend': 'New friends',
        'index.addFriend': 'Add contact',
        'index.search': 'Search',
        'index.searchUserPlaceholder': 'Account or nickname',
        'index.greetingPlaceholder': 'Verification message (optional)',
        'index.searching': 'Searching…',
        'index.keywordRequired': 'Please enter an account or nickname',
        'index.emptySearchResult': 'No matching user',
        'index.emptyFriendRequest': 'No new friend requests',
        'index.add': 'Add',
        'index.accept': 'Accept',
        'index.reject': 'Reject',
        'index.relationFriend': 'Already contacts',
        'index.relationRequested': 'Requested',
        'index.relationIncoming': 'Waiting for you',
        'index.relationBlocked': 'Blocked',
        'index.navBlacklist': 'Blacklist',
        'index.searchBlacklist': 'Search blacklist',
        'index.emptyBlacklist': 'Your blacklist is empty',
        'index.blacklistAdd': 'Block',
        'index.blacklistRemove': 'Unblock',
        'index.blacklistConfirm': 'Block {nickname}? Messages from them will not reach you.',
        'index.blacklistRemoveConfirm': 'Unblock {nickname}?',
        'index.messageBlocked': '{nickname} has blocked you, your message was not delivered',
        'index.friendDelete': 'Delete',
        'index.friendDeleteConfirm': 'Delete {nickname}? Your conversation and chat history will be cleared (their side is not affected); adding them again starts an empty conversation.',
        'index.messageNotFriend': 'You are not contacts, your message was not delivered',
        'index.messageUnfriended': 'This user removed you from their contacts, your message was not delivered',
        'index.defaultGreeting': 'would like to add you as a contact',
        'index.newFriendRequest': '{nickname} would like to add you as a contact',
        'index.friendAccepted': '{nickname} accepted your friend request',
        'index.emojiMessage': '[Sticker]',
        'index.emoji': 'Emoji',
        'index.send': 'Send',
        'index.sendTitle': 'ctrl+Enter to send',
        'index.selectTarget': 'Please select a contact first',
        'index.emptyContent': 'Please enter a message',
        'index.warmTip': 'Notice',
        'index.socketUnsupported': 'Your browser does not support socket'
    }
};

/**
 * 后端接口返回的中文提示（StatusCode != 0 时的 Message）映射。
 * key 为后端原文，未收录的提示原样展示。
 */
I18nUtil.SERVER_MESSAGES = {
    '请求成功': {'en-US': 'Success'},
    '致命错误': {'en-US': 'Fatal error'},
    '暂无权限，请先登录': {'en-US': 'Unauthorized, please sign in first'},
    '账号在别处登录，尝试重新登录': {'en-US': 'Your account was signed in elsewhere, please sign in again'},
    '用户名或密码为空': {'en-US': 'Username or password cannot be empty'},
    '用户名或密码错误': {'en-US': 'Incorrect username or password'},
    '用户名已存在': {'en-US': 'Username already exists'},
    '用户名需为 6-20 位字母、数字、下划线或减号组合': {'en-US': 'Username must be 6-20 letters, digits, underscore or hyphen'},
    '密码至少 6 位': {'en-US': 'Password must be at least 6 characters'},
    '登录成功': {'en-US': 'Signed in successfully'},
    '注册成功': {'en-US': 'Signed up successfully'},
    '退出登录成功': {'en-US': 'Signed out successfully'},
    '登录失败，请稍后重试': {'en-US': 'Sign in failed, please try again later'},
    '注册失败，请稍后重试': {'en-US': 'Sign up failed, please try again later'},
    '退出登录失败': {'en-US': 'Failed to sign out'},
    'token 续签失败': {'en-US': 'Failed to refresh token'},
    '参数错误': {'en-US': 'Invalid parameters'},
    '获取好友列表失败': {'en-US': 'Failed to load contacts'},
    '获取会话列表失败': {'en-US': 'Failed to load conversations'},
    '获取聊天记录失败': {'en-US': 'Failed to load chat history'},
    '标记已读成功': {'en-US': 'Marked as read'},
    '标记已读失败': {'en-US': 'Failed to mark as read'},
    '搜索用户失败': {'en-US': 'Failed to search users'},
    '获取好友申请失败': {'en-US': 'Failed to load friend requests'},
    '好友申请发送失败，请稍后重试': {'en-US': 'Failed to send the friend request, please try again later'},
    '好友申请已发送，等待对方同意': {'en-US': 'Friend request sent, waiting for approval'},
    '你们已成为好友': {'en-US': 'You are contacts now'},
    '已同意好友申请': {'en-US': 'Friend request accepted'},
    '已拒绝好友申请': {'en-US': 'Friend request rejected'},
    '操作失败，请稍后重试': {'en-US': 'Operation failed, please try again later'},
    '用户不存在': {'en-US': 'User does not exist'},
    '不能添加自己为好友': {'en-US': 'You cannot add yourself'},
    '你们已经是好友了': {'en-US': 'You are already contacts'},
    '好友申请不存在': {'en-US': 'Friend request not found'},
    '该好友申请已处理': {'en-US': 'This friend request has been handled'},
    '获取黑名单失败': {'en-US': 'Failed to load the blacklist'},
    '已拉入黑名单': {'en-US': 'Blocked'},
    '拉入黑名单失败，请稍后重试': {'en-US': 'Failed to block the user, please try again later'},
    '不能将自己拉入黑名单': {'en-US': 'You cannot block yourself'},
    '已移出黑名单': {'en-US': 'Unblocked'},
    '移出黑名单失败，请稍后重试': {'en-US': 'Failed to unblock the user, please try again later'},
    '该用户不在黑名单中': {'en-US': 'This user is not in your blacklist'},
    '请先将对方移出黑名单': {'en-US': 'Please unblock this user first'},
    '对方拒绝接收好友申请': {'en-US': 'The user rejected your friend request'},
    '已删除好友': {'en-US': 'Contact deleted'},
    '删除好友失败，请稍后重试': {'en-US': 'Failed to delete the contact, please try again later'},
    '你们还不是好友': {'en-US': 'You are not contacts'}
};

// 语言标识归一化：zh* -> zh-CN，en* -> en-US，其余回落默认语言
I18nUtil.normalize = function (lang) {
    if (!lang) {
        return I18nUtil.DEFAULT_LANG;
    }
    var v = String(lang).toLowerCase();
    if (v.indexOf('zh') === 0) {
        return 'zh-CN';
    }
    if (v.indexOf('en') === 0) {
        return 'en-US';
    }
    return I18nUtil.DEFAULT_LANG;
};

// 文本占位符替换：{name}
I18nUtil.format = function (text, params) {
    if (!params) {
        return text;
    }
    return String(text).replace(/\{(\w+)\}/g, function (match, key) {
        return (params[key] === undefined || params[key] === null) ? match : params[key];
    });
};

// 未存储时的语言探测：localStorage -> 浏览器语言 -> 默认语言
I18nUtil.prototype.detect = function () {
    var v = null;
    try {
        v = localStorage.getItem(I18nUtil.STORAGE_KEY);
    } catch (e) {
        v = null;
    }
    if (v) {
        return v;
    }
    return navigator.language || navigator.userLanguage || I18nUtil.DEFAULT_LANG;
};

I18nUtil.prototype.get = function () {
    return I18nUtil.normalize(this.detect());
};

I18nUtil.prototype.set = function (lang) {
    this.lang = I18nUtil.normalize(lang);
    try {
        localStorage.setItem(I18nUtil.STORAGE_KEY, this.lang);
    } catch (e) {
        console.log(e);
    }
    this.applyHtmlLang(this.lang);
    return this.lang;
};

// 取文案：lang 不传则用实例语言
I18nUtil.prototype.t = function (key, params, lang) {
    var target = I18nUtil.normalize(lang || this.lang || this.get());
    var dict = I18nUtil.MESSAGES[target] || {};
    var text = dict[key];
    if (text === undefined) {
        text = I18nUtil.MESSAGES[I18nUtil.DEFAULT_LANG][key];
    }
    if (text === undefined) {
        return key;
    }
    return I18nUtil.format(text, params);
};

// 翻译后端返回的提示，未收录则原样返回
I18nUtil.prototype.tServer = function (msg, lang) {
    if (!msg) {
        return '';
    }
    var target = I18nUtil.normalize(lang || this.lang || this.get());
    var dict = I18nUtil.SERVER_MESSAGES[msg];
    if (!dict || !dict[target]) {
        return msg;
    }
    return dict[target];
};

I18nUtil.prototype.htmlLang = function (lang) {
    return I18nUtil.HTML_LANGS[I18nUtil.normalize(lang || this.lang)] || 'zh-CN';
};

I18nUtil.prototype.applyHtmlLang = function (lang) {
    document.documentElement.setAttribute('lang', this.htmlLang(lang));
};
