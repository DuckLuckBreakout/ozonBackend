import AjaxModule from '../modules/Ajax/Ajax.js';
import {serverApiPath, urls} from '../utils/urls/urls';

import BaseModel from './BaseModel';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import HTTPResponses from '../utils/http-responses/httpResponses';
import Bus from '../utils/bus/bus';

/**
 * @description Model for Signup in MVP Arch
 */
class SignupModel extends BaseModel {
    /**
     *
     * @param {string} email
     * @param {string} password
     */
    signupUser({email, password}) {
        AjaxModule.postUsingFetch({
            url: serverApiPath + urls.signupUrl,
            body: {email, password},
        }).then((response) => {
            switch (response.status) {
            case HTTPResponses.Created: {
                Bus.globalBus.emit(Events.SignupEmitResult, Responses.Success);
                this.bus.emit(Events.SignupEmitResult, Responses.Success);
                break;
            }
            case HTTPResponses.Offline: {
                this.bus.emit(Events.SignupEmitResult, Responses.Offline);
                break;
            }
            case HTTPResponses.Conflict: {
                this.bus.emit(Events.SignupEmitResult, Responses.Conflict);
                break;
            }
            default: {
                throw Responses.Error;
            }
            }
        }).catch(() => {
            this.bus.emit(Events.SignupEmitResult, Responses.Error);
        });
    }
}

export default SignupModel;
