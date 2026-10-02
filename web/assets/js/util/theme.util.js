function ThemeUtil() {
}

ThemeUtil.STORAGE_KEY = 'mini-im-theme';

ThemeUtil.prototype.get = function () {
    var v = localStorage.getItem(ThemeUtil.STORAGE_KEY)
    return (v === 'light' || v === 'dark' || v === 'system') ? v : 'system'
}

ThemeUtil.prototype.set = function (v) {
    if (v !== 'light' && v !== 'dark' && v !== 'system') {
        v = 'system'
    }
    localStorage.setItem(ThemeUtil.STORAGE_KEY, v)
    this.apply()
}

ThemeUtil.prototype.resolve = function (v) {
    if (v === 'system') {
        v = (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) ? 'dark' : 'light'
    }
    return v
}

ThemeUtil.prototype.apply = function () {
    document.documentElement.setAttribute('data-theme', this.resolve(this.get()))
}

ThemeUtil.prototype.init = function () {
    var self = this
    this.apply()
    if (window.matchMedia) {
        var mq = window.matchMedia('(prefers-color-scheme: dark)')
        var handler = function () {
            if (self.get() === 'system') {
                self.apply()
            }
        }
        if (mq.addEventListener) {
            mq.addEventListener('change', handler)
        } else if (mq.addListener) {
            mq.addListener(handler)
        }
    }
}
