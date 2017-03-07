import 'material-design-lite'

function applyTemplate (selector, html, append=false) {
    if (append) {
        $(selector).append(html)
    } else {
        $(selector).html(html)
    }
    componentHandler.upgradeElements($(selector).children())
}

export {
    applyTemplate
}