class CryptUtil {
    encrypt(value) {
        return window.btoa(window.encodeURIComponent(JSON.stringify(value)))
    }

    decrypt(value) {
        return JSON.parse(window.decodeURIComponent(window.atob(value)))
    }
}
