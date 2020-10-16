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
}

export default service