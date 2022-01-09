import FLNameModule from './subModules/FLNameModule';
import AvatarModule from './subModules/AvatarModule';
import EmailModule from './subModules/EmailModule';
import CheckAuthModule from './subModules/CheckAuthModule';
import BaseModel from './BaseModel';
import AjaxModule from '../modules/Ajax/Ajax';
import {serverApiPath, urls} from '../utils/urls/urls';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import HTTPResponses from '../utils/http-responses/httpResponses';
import Bus from '../utils/bus/bus.js';

/**
 * @description Model for Profile in MVP Arch. THIS IS A FACADE!
 */
class ProfileModel extends BaseModel {
    #fLNameModel;
    #emailModel;
    #avatarModel;
    #checkAuthModel;

    /**
     * @param {Object} bus bus of this mvp part
     * @description simple Constructor. Creating all elements inside (FACADE)
     */
    constructor(bus) {
        super(bus);
        this.#fLNameModel = new FLNameModule(bus);
        this.#emailModel = new EmailModule(bus);
        this.#avatarModel = new AvatarModule(bus);
        this.#checkAuthModel = new CheckAuthModule(bus);
    }

    /**
     * @return {string}
     */
    get firstName() {
        return this.#fLNameModel.firstName;
    }

    /**
     *
     * @return {string}
     */
    get lastName() {
        return this.#fLNameModel.lastName;
    }

    /**
     *
     * @return {string|*}
     */
    get email() {
        return this.#emailModel.email;
    }

    /**
     *
     * @return {string}
     */
    get avatarURL() {
        return this.#avatarModel.avatarURL;
    }

    /**
     *
     * @param {string} firstName
     * @param {string} lastName
     */
    changeFirstLastName = (firstName, lastName) => {
        this.#fLNameModel.changeFirstLastName(firstName, lastName);
    }

    /**
     * @return {void}
     * @description We will get Name later by emitiing a signal!
     */
    getFirstLastName = () => {
        this.#fLNameModel.getFirstLastName();
    }

    /**
     *
     * @param {File} file
     */
    changeAvatar = (file) => {
        this.#avatarModel.changeAvatar(file);
    }

    /**
     * @return {string}
     */
    getAvatar = () => {
        return this.avatarURL;
    }

    /**
     * @description Checking for Auth. Emit signal later
     */
    checkAuth = () => {
        this.#checkAuthModel.checkAuth();
    }

    /**
     * @description Clears all saved data in model
     */
    clear() {
        this.#fLNameModel.clear();
        this.#emailModel.clear();
        this.#avatarModel.clear();
        this.#checkAuthModel.clear();
    }

    /**
     * @description gets all data for profile from one AJAX request
     */
    getProfileData() {
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.profileUrl,
            body: null,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((response) => {
            window.localStorage.setItem('firstName', response.first_name);
            this.#fLNameModel.firstName = response.first_name;
            window.localStorage.setItem('lastName', response.last_name);
            this.#fLNameModel.lastName = response.last_name;
            window.localStorage.setItem('avatarURL', response.avatar.url);
            this.#avatarModel.avatarURL = response.avatar.url;
            window.localStorage.setItem('email', response.email);
            this.#emailModel.email = response.email;
            this.bus.emit(Events.ProfileAllResult, Responses.Success);
        }).catch((result) => {
            switch (result) {
            case Responses.Offline: {
                this.bus.emit(Events.ProfileAllResult, Responses.Offline);
                break;
            }
            default: {
                this.bus.emit(Events.ProfileAllResult, Responses.Error);
                break;
            }
            }
        });
    }

    /**
     * @description logout profile from one AJAX request
     */
    profileLogout() {
        AjaxModule.deleteUsingFetch({
            url: serverApiPath + urls.logoutUrl,
        }).then((response) => {
            if (response.status === HTTPResponses.Success) {
                Bus.globalBus.emit(Events.ProfileLogoutEmitResult, Responses.Success);
            } else {
                Bus.globalBus.emit(Events.ProfileLogoutEmitResult, Responses.Error);
            }
        }).catch(() => {
            Bus.globalBus.emit(Events.ProfileLogoutEmitResult, Responses.Error);
        });
    }
}

export default ProfileModel;
