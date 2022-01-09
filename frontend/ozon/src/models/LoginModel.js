import AjaxModule from '../modules/Ajax/Ajax.js';
import {serverApiPath, urls} from '../utils/urls/urls';
import BaseModel from './BaseModel';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import HTTPResponses from '../utils/http-responses/httpResponses';
import Bus from '../utils/bus/bus';

/**
 * @description Model for Log in User in MVP Arch
 */
class LoginModel extends BaseModel {
    /**
     *
     * @param {string} email
     * @param {string} password
     */
    loginUser({email, password}) {
        AjaxModule.postUsingFetch({
            url: serverApiPath + urls.loginUrl,
            body: {email, password},
        }).then((response) => {
            switch (response.status) {
            case HTTPResponses.Success: {
                Bus.globalBus.emit(Events.LoginEmitResult, Responses.Success);
                this.bus.emit(Events.LoginEmitResult, Responses.Success);
                break;
            }
            case HTTPResponses.Offline: {
                this.bus.emit(Events.LoginEmitResult, Responses.Offline);
                break;
            }
            case HTTPResponses.Unauthorized: {
                this.bus.emit(Events.LoginEmitResult, Responses.Unauthorized);
                break;
            }
            default: {
                throw Responses.Error;
            }
            }
        }).catch(() => {
            this.bus.emit(Events.LoginEmitResult, Responses.Error);
        });
    }
}

export default LoginModel;
