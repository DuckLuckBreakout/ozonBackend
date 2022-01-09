import BasePresenter from './BasePresenter.js';
import {isValidInputs} from '../modules/Valiadtor/validator';
import Router from '../utils/router/Router';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import Bus from '../utils/bus/bus';

/**
 * @description Presenter for Profile View and Model
 */
class ProfilePresenter extends BasePresenter {
    /**
     * @param {HTMLElement} application html of application
     * @param {Class} View Class of view object
     * @param {Class} Model Class of model object
     */
    constructor(application, View, Model) {
        super(application, View, Model);
        this.bus.on(Events.ProfileFLNameChange, this.sendFirstLastName);
        this.bus.on(Events.ProfileFLNameResult, this.firstLastNameSendProcessResult);
        this.bus.on(Events.ProfileAvatarChange, this.sendAvatar);
        this.bus.on(Events.ProfileAvatarResult, this.avatarSendProcessResult);
        this.bus.on(Events.ProfileEmailResult, this.emailSendProcessResult);
        this.bus.on(Events.ProfileCheckAuthResult, this.tryAuthProcessResult);
        this.bus.on(Events.ProfileAllGet, this.getAllData);
        this.bus.on(Events.ProfileAllResult, this.renderAllData);
        this.bus.on(Events.ProfileLogoutPrepare, this.profileLogoutPrepare);

        Bus.globalBus.on(Events.ProfileNewUserLoggedIn, this.removeData);
        Bus.globalBus.on(Events.ProfileTransmitData, this.returnUserData);
        Bus.globalBus.on(Events.ProfileLogout, this.profileLogout);
    }

    /**
     * @param {Object} specificTypeToCheck
     * @return {Object} contains array of invalid fields
     */
    isFirstLasNameFormValid = (specificTypeToCheck = []) => {
        return isValidInputs(document.getElementById('flname-form').getElementsByTagName('input'),
            specificTypeToCheck);
    }

    /**
     * @return {boolean} is dispatched file valid?
     */
    isAvatarValid = () => {
        const input = document.createElement('input');
        input.setAttribute('type', 'file');
        input.files = document.getElementById('avatar-input').files;
        return !!isValidInputs([input]).failedFields.length;
    }

    /**
     * @description Send data to model
     */
    sendFirstLastName = () => {
        const result = this.isFirstLasNameFormValid();
        if (result.failedFields.length) {
            this.view.invalidForm(result.failedFields);
            return;
        }
        const firstName = document.getElementsByName('firstName')[0].value.trim();
        const lastName = document.getElementsByName('lastName')[0].value.trim();
        this.model.changeFirstLastName(firstName, lastName);
    }

    /**
     * @param {string} result
     * @description Processed a signal
     */
    firstLastNameSendProcessResult = (result) => {
        switch (result) {
        case Responses.Success: {
            this.view.changeFirstLastName(this.getFirstName(), this.getLastName());
            break;
        }
        case Responses.Offline: {
            Router.open('/offline');
            break;
        }
        default: {
            console.error(result);
            break;
        }
        }
    }

    /**
     *
     * @return {string}
     */
    getFirstName = () => {
        if (this.model.firstName !== undefined) {
            return this.model.firstName;
        }

        this.model.getFirstLastName();
        return '';
    }

    /**
     *
     * @return {string}
     */
    getLastName = () => {
        if (this.model.lastName !== undefined) {
            return this.model.lastName;
        }

        this.model.getFirstLastName();
        return '';
    }

    /**
     *
     * @return {string}
     */
    getEmail = () => {
        if (this.model.email !== undefined) {
            return this.model.email;
        }

        this.model.getEmail();
        return '';
    }

    /**
     *
     * @return {string}
     */
    getAvatar = () => {
        if (this.model.avatarURL !== undefined) {
            return this.model.avatarURL;
        }

        this.model.getAvatar();
        return '';
    }

    /**
     * @description get data view and send to model
     */
    sendAvatar = () => {
        if (this.isAvatarValid()) {
            this.view.invalidAvatar();
            return;
        }
        const avatarInput = document.getElementById('avatar-input');
        this.model.changeAvatar(avatarInput.files[0]);
    }

    /**
     * @description logout from profile
     */
    profileLogoutPrepare = () => {
        if (!navigator.onLine) {
            Router.open('/offline');
            return;
        }
        this.view.remove();
        Bus.globalBus.emit(Events.HeaderSetCartItems, 0);
        Bus.globalBus.emit(Events.CartDrop);
        Bus.globalBus.emit(Events.WebPushUnsubscribe);
    }

    profileLogout = () => {
        this.model.profileLogout();
        Router.open('/', {replaceState: true});
    }


    /**
     *
     * @param {string} result
     */
    avatarSendProcessResult = (result) => {
        switch (result) {
        case Responses.Success: {
            this.view.changeAvatar(this.getAvatar());
            break;
        }
        case Responses.Offline: {
            Router.open('/offline');
            break;
        }
        default: {
            console.error(result);
            break;
        }
        }
    }

    /**
     *
     * @param {string} result
     */
    emailSendProcessResult = (result) => {
        switch (result) {
        case Responses.Success: {
            this.view.changeEmail(this.getEmail());

            break;
        }
        case Responses.Offline: {
            Router.open('/offline');
            break;
        }
        default: {
            console.error(result);
            break;
        }
        }
    }

    /**
     * @description attempts to authorize
     */
    tryAuth = () => {
        this.model.checkAuth();
    }

    /**
     *
     * @param {string} result
     */
    tryAuthProcessResult = (result) => {
        switch (result) {
        case Responses.Success: {
            this.view.render();
            this.view.cache.hidden = false;
            break;
        }
        case Responses.Offline: {
            Router.open('/offline');
            break;
        }
        default: {
            Router.open('/login', {replaceState: true});
            break;
        }
        }
    }

    /**
     * @description This function is used when user is re-loginned or creates a new account
     * Just clear cache and model
     */
    removeData = () => {
        this.model.clear();
        this.view.clearCache();
    }

    getAllData = () => {
        this.model.getProfileData();
    }

    /**
     *
     * @param {Responses} status
     */
    renderAllData = (status) => {
        if (status === Responses.Success) {
            this.view.renderData();
            return;
        }
        if (status === Responses.Offline) {
            Router.open('/offline', {replaceState: true});
            return;
        }
        console.error(status);
    }

    /**
     * @param {string} eventToEmit
     */
    returnUserData = (eventToEmit) => {
        Bus.globalBus.emit(eventToEmit, {
            firstName: this.model.firstName,
            lastName: this.model.lastName,
            email: this.model.email,
            avatarURL: this.model.avatarURL,
        });
    }
}

export default ProfilePresenter;
