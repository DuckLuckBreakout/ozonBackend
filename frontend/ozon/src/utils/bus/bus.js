/**
 * @description Simple Bus class which provider "Mediator" pattern through all application.
 */
class Bus {
    /**
     * @description Simple constructor
     */
    static #globalBus


    /**
     * @return {Object} Singleton. Global Bus to cross-presenter actions
     */
    static get globalBus() {
        if (Bus.#globalBus === undefined) {
            Bus.#globalBus = new Bus();
        }

        return Bus.#globalBus;
    }

    /**
     * @description default constructor
     */
    constructor() {
        this.listeners = {};
    }

    /**
     *
     * @param {string} event Some event
     * @param {function} callback A function which will call when event is emitted.
     */
    on(event, callback) {
        this.listeners[event] = this.listeners[event] || [];
        this.listeners[event].push(callback);
    }

    /**
     *
     * @param {string} event Some event
     * @param {function} callback A function which will removed from callable when event is emitted.
     */
    off(event, callback) {
        this.listeners[event] = this.listeners[event]
            .filter((listener) => {
                return listener !== callback;
            });
    }

    /**
     *
     * @param {string} event Some event
     * @param {Object} data Data which will be passed through function.
     */
    emit(event, ...data) {
        this.listeners[event].forEach((listener) => {
            listener(...data);
        });
    }
}

export default Bus;
