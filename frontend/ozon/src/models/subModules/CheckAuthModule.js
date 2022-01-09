import AjaxModule from '../../modules/Ajax/Ajax';
import {serverApiPath, urls} from '../../utils/urls/urls';
import BaseModel from '../BaseModel';
import Events from '../../utils/bus/events';
import Responses from '../../utils/bus/responses';
import HTTPResponses from '../../utils/http-responses/httpResponses';

/**
 * @description Model for Checking for Users Authorization in MVP Arch
 */
class CheckAuthModule extends BaseModel {
    /**
     * @description Sending Ajax. Then emits signal of success/failure of Authorization
     */
    checkAuth() {
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.profileUrl,
            body: null,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw Responses.Error;
            }
            this.bus.emit(Events.ProfileCheckAuthResult, Responses.Success);
        }).catch(() => {
            this.bus.emit(Events.ProfileCheckAuthResult, Responses.Error);
        });
    }
}

export default CheckAuthModule;
