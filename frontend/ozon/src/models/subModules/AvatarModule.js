import AjaxModule from '../../modules/Ajax/Ajax.js';
import BaseModel from '../BaseModel';
import Events from '../../utils/bus/events';
import Responses from '../../utils/bus/responses';
import HTTPResponses from '../../utils/http-responses/httpResponses';
import {serverApiPath, urls} from '../../utils/urls/urls';

/**
 * @description Model for Avatar Loading/Uploading in MVP Arch
 */
class AvatarModule extends BaseModel {
    #avatarURL

    /**
     *
     * @param {string} URL New URL for avatar
     */
    set avatarURL(URL) {
        if (!URL.includes('avatar')) {
            this.#avatarURL = urls.defaultAvatar;
            return;
        }

        window.localStorage.setItem('avatarURL', URL);
        this.#avatarURL= URL;
    }

    /**
     *
     * @return {string} Avatar`s URL
     */
    get avatarURL() {
        if (this.#avatarURL !== undefined) {
            return this.#avatarURL;
        }
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.profileAvatarUrl,
            body: null,
        }).then((response) => {
            return response.json();
        }).then((response) => {
            this.#saveAndEmit(response.url);
        }).catch(() => {
            this.bus.emit(Events.ProfileAvatarResult, Responses.Error);
        });
        return urls.defaultAvatar;
    }

    /**
     *
     * @param {File} file A file which will uploaded to server
     */
    changeAvatar = (file) => {
        const formData = new FormData();
        formData.append('avatar', file);
        AjaxModule.putUsingFetch({
            data: true,
            url: serverApiPath + urls.profileAvatarUrl,
            body: formData,
        }).then(() => {
            AjaxModule.getUsingFetch({
                url: serverApiPath + urls.profileAvatarUrl,
                body: null,
            }).then((response) => {
                if (response.status !== HTTPResponses.Success) {
                    throw response.status;
                }
                return response.json();
            }).then((response) => {
                this.#saveAndEmit(response.url);
            });
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Unauthorized: {
                this.bus.emit(Events.ProfileAvatarResult, Responses.Unauthorized);
                break;
            }
            case HTTPResponses.Offline: {
                this.bus.emit(Events.ProfileAvatarResult, Responses.Offline);
                break;
            }
            default: {
                this.bus.emit(Events.ProfileAvatarResult, Responses.Error);
                break;
            }
            }
        });
    }

    /**
     *
     * @param {string} url url to avatar on file server
     */
    #saveAndEmit = (url) => {
        this.avatarURL = url;
        this.bus.emit(Events.ProfileAvatarResult, Responses.Success);
    }

    /**
     * @description Clears all saved data in model
     */
    clear() {
        this.#avatarURL = undefined;
    }
}

export default AvatarModule;

