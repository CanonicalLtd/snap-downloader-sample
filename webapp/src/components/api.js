import axios from 'axios'
import constants from './constants'

let service = {
    authGet: () => {
        return axios.get(constants.baseUrl + 'auth');
    },

    storeLogin: (email, password, otp, storeId, series) => {
        return axios.post(constants.baseUrl + 'login', {
            email: email, password: password, otp: otp, store: storeId, series: series});
    },

    snapsList: () => {
        return axios.get(constants.baseUrl + 'snaps');
    },

    snapCreate: (name, arch) => {
        return axios.post(constants.baseUrl + 'snaps', {name: name, arch: arch});
    },

    snapDelete: (id) => {
        return axios.delete(constants.baseUrl + 'snaps/' + id);
    },

    snapDownload: (name, filename) => {
        return axios.get(constants.baseUrl + 'downloads/' + name + '/' + filename);
    },

    snapsDownloadList: () => {
        return axios.get(constants.baseUrl + 'downloads');
    },

    settingsInterval: () => {
        return axios.get(constants.baseUrl + 'settings/interval');
    },

    settingsSetInterval: (interval) => {
        return axios.post(constants.baseUrl + 'settings/interval', {value: interval});
    },

    settingsLastRun: () => {
        return axios.get(constants.baseUrl + 'settings/lastrun');
    },
}

export default service