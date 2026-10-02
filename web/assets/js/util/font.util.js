/**
 * 字号 / 字重设置工具
 *
 * 与 theme.util.js 同构：取值存 localStorage，并把结果落到 <html> 上：
 *   字号 -> 直接写根字号内联样式（无单位字号的 rem 基准），同时在 data-font-size 上留个记号
 *   字重 -> data-font-weight，切 common.css 里的 --fw-base / --fw-medium / --fw-bold
 *
 * 字号是 10~30 之间的整数（px），抽屉里用滑杆连续调节；CSS 里的字号都是 rem，
 * 所以改根字号只缩放文字，间距与图片尺寸不变。
 *
 * 用法：
 *   var fontUtil = new FontUtil()
 *   fontUtil.init()                 // 进页面时调用一次
 *   fontUtil.getSize()              // 14
 *   fontUtil.setSize('18')          // 保存并立即生效，返回归一化后的 18
 *   fontUtil.getWeight()            // 'normal' | 'bold'
 *   fontUtil.setWeight('bold')
 */
function FontUtil() {
}

// 存储键 / 属性名要和页面 <head> 里的内联脚本保持一致（避免首屏闪一下默认字号）
FontUtil.SIZE_KEY = 'mini-im-font-size';
FontUtil.WEIGHT_KEY = 'mini-im-font-weight';
FontUtil.SIZE_ATTR = 'data-font-size';
FontUtil.WEIGHT_ATTR = 'data-font-weight';

FontUtil.MIN_SIZE = 10;
FontUtil.MAX_SIZE = 30;
FontUtil.DEFAULT_SIZE = 14;

FontUtil.WEIGHTS = ['normal', 'bold'];
FontUtil.DEFAULT_WEIGHT = 'normal';

// 字号归一化：取整 + 夹到 10~30，非法值回落默认
FontUtil.normalizeSize = function (value) {
    var size = parseInt(value, 10)
    if (isNaN(size)) {
        size = FontUtil.DEFAULT_SIZE
    }
    return Math.min(FontUtil.MAX_SIZE, Math.max(FontUtil.MIN_SIZE, size))
};

FontUtil.normalizeWeight = function (value) {
    return FontUtil.WEIGHTS.indexOf(value) > -1 ? value : FontUtil.DEFAULT_WEIGHT;
};

FontUtil.prototype.getSize = function () {
    return FontUtil.normalizeSize(localStorage.getItem(FontUtil.SIZE_KEY))
}

FontUtil.prototype.setSize = function (value) {
    var size = FontUtil.normalizeSize(value)
    localStorage.setItem(FontUtil.SIZE_KEY, size)
    this.apply()
    return size
}

FontUtil.prototype.getWeight = function () {
    return FontUtil.normalizeWeight(localStorage.getItem(FontUtil.WEIGHT_KEY))
}

FontUtil.prototype.setWeight = function (value) {
    var weight = FontUtil.normalizeWeight(value)
    localStorage.setItem(FontUtil.WEIGHT_KEY, weight)
    this.apply()
    return weight
}

FontUtil.prototype.apply = function () {
    var size = this.getSize()
    document.documentElement.style.fontSize = size + 'px'
    document.documentElement.setAttribute(FontUtil.SIZE_ATTR, size)
    document.documentElement.setAttribute(FontUtil.WEIGHT_ATTR, this.getWeight())
}

FontUtil.prototype.init = function () {
    this.apply()
}
