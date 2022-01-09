/**
 * @description Base Model in MVP Arch
 */
class BaseModel {
    #bus

    /**
     * @param {Object} bus bus of this mvp part
     */
    constructor(bus) {
        this.#bus = bus;
    }

    /**
     * @return {Object} bus of this mvp part
     */
    get bus() {
        return this.#bus;
    }

    /**
     * @description Clears all saved data in model
     */
    clear() {

    }
}

export default BaseModel;
