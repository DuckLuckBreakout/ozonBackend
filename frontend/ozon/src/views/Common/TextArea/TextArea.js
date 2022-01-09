/**
 * @class TextArea
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class TextArea {
    /**
     *
     * @param {string} name name of a text area
     * @param {number} rows amount of rows in text area
     * @param {number} cols amount of columns in text aread
     */
    constructor(
        {
            name,
            rows,
            cols,
        } = {}) {
        this.objectType = 'textarea';
        this.name = name;
        this.rows = rows;
        this.cols = cols;
    }
}

export default TextArea;
