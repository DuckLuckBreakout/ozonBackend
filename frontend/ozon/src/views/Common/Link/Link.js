/**
 * @class Link
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class Link {
    /**
     *
     * @param {string} href link where to go
     * @param {string} name name of a link
     * @param {string} type type of a link
     * @param {string} dataSection
     */
    constructor({href, name, type, dataSection} = {}) {
        this.objectType = 'link';
        this.href = href;
        this.name = name;
        this.dataSection = dataSection;
        this.type = type;
    }
}

export default Link;
