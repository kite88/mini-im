class CacheData {
    mKey = 'messageList'
    fKey = 'friendList'
    cKey = 'chatList_'

    saveMessage(list) {
        localStorage.setItem(this.mKey, JSON.stringify(list))
    }

    getMessage() {
        var data = localStorage.getItem(this.mKey)
        var list = []
        if (data !== null) {
            list = JSON.parse(data)
        }
        return list
    }

    clearMessage() {
        localStorage.removeItem(this.mKey)
    }

    saveFriend(list) {
        localStorage.setItem(this.fKey, JSON.stringify(list))
    }

    getFriend() {
        var data = localStorage.getItem(this.fKey)
        var list = []
        if (data !== null) {
            list = JSON.parse(data)
        }
        return list
    }

    clearFriend() {
        localStorage.removeItem(this.fKey)
    }

    saveChat(key, list) {
        sessionStorage.setItem(this.cKey + key, JSON.stringify(list))
    }

    getChat(key) {
        var data = sessionStorage.getItem(this.cKey + key)
        var list = []
        if (data !== null) {
            list = JSON.parse(data)
        }
        return list
    }

    isChat(key) {
        var data = sessionStorage.getItem(this.cKey + key)
        return data !== null;
    }

    delChat(key) {
        sessionStorage.removeItem(this.cKey + key)
    }
}
