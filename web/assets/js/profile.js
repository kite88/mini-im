class Profile {
    key = 'userInfo';

    save(profile) {
        sessionStorage.setItem(this.key, JSON.stringify(profile))
    }

    get() {
        return JSON.parse(sessionStorage.getItem(this.key))
    }

    del() {
        sessionStorage.removeItem(this.key)
    }

    getUser() {
        var data = this.get()
        if (data) {
            return data.user
        }
        return null
    }

    getToken() {
        var data = this.get()
        if (data) {
            var nowUnixTime = Math.round(new Date() / 1000)
            if (data.tokens.exp - nowUnixTime < 43200) {
                this.refreshToken(data.tokens.token)
            }
            return data.tokens.token
        }
        return ""
    }

    refreshToken(token) {
        axios({
            url: apiDomain + '/api/token',
            method: 'put',
            headers: {
                authorization: 'Bearer ' + token
            },
        }).catch(function (err) {
            console.log(err)
        }).then(res => {
            var tokens = res.data.Data
            var user = this.get()
            user.tokens = tokens
            this.save(user)
        })
    }
}
