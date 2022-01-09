/**
 * @param {number} from start of iteration
 * @param {number} to end of iteration
 * @param {number} incr step of iteration
 * @param {Object} block block from which we gain information
 * @return {string} accumulated value
 */
const forHelper = function(from, to, incr, block) {
    let accum = '';
    for (let i = from; i < to; i += incr) {
        accum += block.fn(i);
    }
    return accum;
};

export default forHelper;
