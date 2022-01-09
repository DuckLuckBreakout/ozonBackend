
/**
 * @class Img
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class Img {
    /**
     * @param {string} URL source of an image
     */
    constructor({src = '',
    } = {}) {
        this.objectType = 'img';
        this.src = src;
    }
}

export default Img;
