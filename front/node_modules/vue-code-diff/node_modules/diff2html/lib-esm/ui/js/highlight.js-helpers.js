function escapeHTML(value) {
    return value.replace(/&/gm, '&amp;').replace(/</gm, '&lt;').replace(/>/gm, '&gt;');
}
function tag(node) {
    return node.nodeName.toLowerCase();
}
export function nodeStream(node) {
    var result = [];
    var nodeStream = function (node, offset) {
        for (var child = node.firstChild; child; child = child.nextSibling) {
            if (child.nodeType === 3 && child.nodeValue !== null) {
                offset += child.nodeValue.length;
            }
            else if (child.nodeType === 1) {
                result.push({
                    event: 'start',
                    offset: offset,
                    node: child,
                });
                offset = nodeStream(child, offset);
                if (!tag(child).match(/br|hr|img|input/)) {
                    result.push({
                        event: 'stop',
                        offset: offset,
                        node: child,
                    });
                }
            }
        }
        return offset;
    };
    nodeStream(node, 0);
    return result;
}
export function mergeStreams(original, highlighted, value) {
    var processed = 0;
    var result = '';
    var nodeStack = [];
    function isElement(arg) {
        var _a;
        return arg !== null && ((_a = arg) === null || _a === void 0 ? void 0 : _a.attributes) !== undefined;
    }
    function selectStream() {
        if (!original.length || !highlighted.length) {
            return original.length ? original : highlighted;
        }
        if (original[0].offset !== highlighted[0].offset) {
            return original[0].offset < highlighted[0].offset ? original : highlighted;
        }
        return highlighted[0].event === 'start' ? original : highlighted;
    }
    function open(node) {
        if (!isElement(node)) {
            throw new Error('Node is not an Element');
        }
        result += "<" + tag(node) + " " + Array()
            .map.call(node.attributes, function (attr) { return attr.nodeName + "=\"" + escapeHTML(attr.value).replace(/"/g, '&quot;') + "\""; })
            .join(' ') + ">";
    }
    function close(node) {
        result += '</' + tag(node) + '>';
    }
    function render(event) {
        (event.event === 'start' ? open : close)(event.node);
    }
    while (original.length || highlighted.length) {
        var stream = selectStream();
        result += escapeHTML(value.substring(processed, stream[0].offset));
        processed = stream[0].offset;
        if (stream === original) {
            nodeStack.reverse().forEach(close);
            do {
                render(stream.splice(0, 1)[0]);
                stream = selectStream();
            } while (stream === original && stream.length && stream[0].offset === processed);
            nodeStack.reverse().forEach(open);
        }
        else {
            if (stream[0].event === 'start') {
                nodeStack.push(stream[0].node);
            }
            else {
                nodeStack.pop();
            }
            render(stream.splice(0, 1)[0]);
        }
    }
    return result + escapeHTML(value.substr(processed));
}
//# sourceMappingURL=highlight.js-helpers.js.map