class DateUtil {
    formatUnix(timestamp) {
        if (timestamp === '' || timestamp === 0 || timestamp === null || timestamp === undefined) {
            return ''
        }
        var data3 = new Date(timestamp * 1000) // 精确到毫秒的就不用 *1000  了
        var y = data3.getFullYear(),
            month = data3.getMonth() + 1,
            m = month < 10 ? '0' + month : month,
            d = (data3.getDate() < 10) ? '0' + data3.getDate() : data3.getDate(),
            h = (data3.getHours() < 10) ? '0' + data3.getHours() : data3.getHours(),
            mi = (data3.getMinutes() < 10) ? '0' + data3.getMinutes() : data3.getMinutes(),
            s = (data3.getSeconds() < 10) ? '0' + data3.getSeconds() : data3.getSeconds()
        return y + '/' + m + '/' + d + ' ' + h + ':' + mi
    }

    formatUnixDate(timestamp) {
        return this.formatUnix(timestamp) === '' ? '' : this.formatUnix(timestamp).toString().split(' ')[0];
    }

    formatUnixTime(timestamp) {
        return this.formatUnix(timestamp) === '' ? '' : this.formatUnix(timestamp).toString().split(' ')[1];
    }

    formatTime(timestamp) {
        if (timestamp === '' || timestamp === 0 || timestamp === null || timestamp === undefined) {
            return ''
        }
        var n_date = new Date()
        var n_year = n_date.getFullYear(),
            n_month = n_date.getMonth() + 1,
            n_day = n_date.getDate()
        var c_date = new Date(timestamp * 1000)
        var c_year = c_date.getFullYear(),
            c_month = c_date.getMonth() + 1,
            c_day = c_date.getDate(),
            c_hours = c_date.getHours(),
            c_minutes = c_date.getMinutes()
        c_minutes = c_minutes < 10 ? '0' + c_minutes : c_minutes
        if ((n_year.toString() + n_month.toString() + n_day.toString()) === (c_year.toString() + c_month.toString() + c_day.toString())) {
            return c_hours + ':' + c_minutes
        }
        if (n_year.toString() === c_year.toString()) {
            return c_month + '/' + c_day + ' ' + c_hours + ':' + c_minutes
        }
        return c_year.toString().substring(2, 4) + '/' + c_month + '/' + c_day + ' ' + c_hours + ':' + c_minutes
    }
}
