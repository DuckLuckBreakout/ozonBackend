import AjaxModule from '../modules/Ajax/Ajax';
import {serverApiPath, urls} from '../utils/urls/urls';
import BaseModel from './BaseModel';
import HTTPResponses from '../utils/http-responses/httpResponses';
import urlBase64ToUint8Array from '../utils/urlBase64ToUint8Array/urlBase64ToUint8Array';
import Bus from '../utils/bus/bus';
import Events from '../utils/bus/events';

/**
 * @description Model for WebPush notification from server
 * Attention! This model does not have a presenter. Events from the global bus are handled by it itself
 */
class WebPushModel extends BaseModel {
    #publicKey;

    /**
     * @description simple constructor.
     */
    constructor() {
        super();
        Bus.globalBus.on(Events.WebPushSubscribe, this.register);
        Bus.globalBus.on(Events.WebPushUnsubscribe, this.unsubscribe);
    }

    register = () => {
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.publicKey,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((response) => {
            this.#publicKey = response.key;
            this.registerSendRequest();
        }).catch((err) => {
            console.error(err);
        });
    }

    unsubscribe = () => {
        if (this.#publicKey) {
            this.unsubscribeSendRequest();
            return;
        }
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.publicKey,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((response) => {
            this.#publicKey = response.key;
            this.unsubscribeSendRequest();
        }).catch((err) => {
            console.error(err);
        });
    }

    /**
     * @return {Promise<void>}
     */
    async registerSendRequest() {
        const register = await navigator.serviceWorker.getRegistration();
        const subscription = await register.pushManager.subscribe({
            userVisibleOnly: true,
            applicationServerKey: urlBase64ToUint8Array(this.#publicKey),
        });

        await AjaxModule.postUsingFetch({
            url: serverApiPath + urls.notifications,
            body: subscription,
        });
    }

    /**
     * @return {Promise<void>}
     */
    async unsubscribeSendRequest() {
        const register = await navigator.serviceWorker.getRegistration();
        const subscription = await register.pushManager.subscribe({
            userVisibleOnly: true,
            applicationServerKey: urlBase64ToUint8Array(this.#publicKey),
        });

        await AjaxModule.deleteUsingFetch({
            url: serverApiPath + urls.notifications,
            body: subscription,
        }).then((_) => {
            Bus.globalBus.emit(Events.ProfileLogout);
        }).catch((err) => {
            console.error(err);
            Bus.globalBus.emit(Events.ProfileLogout);
        });
    }
}

export default WebPushModel;
