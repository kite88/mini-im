class CookieUtil {
    set(name, value, second) {
        var life = new Date().getTime() + second * 1000;
        document.cookie = name + "=" + escape(value) + ";expires=" + new Date(life).toGMTString();
    }

    get(name) {
        var cookies = document.cookie.split(";");
        if (cookies.length > 0) {
            var cookie = cookies[0].split("=");
            if (cookie[0] === name)
                return unescape(cookie[1]);
        }
        return null;
    }

    del(name) {
        document.cookie = name + "=" + escape('null') + ";expires=" + new Date().toGMTString();
    }
}
