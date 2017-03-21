import 'styles/main.scss'
import jQuery from 'jquery'
import ccPayment from 'payment/credit-card'

window.$ = jQuery
window._currentLoaderId = undefined
const clearPageContent = () => {$(".page-content").html("")}
const loaders = {
    "sidebarHome": {
        icon: "home",
        label: "Home",
        render: () => {
            $.ajax('/api/hello').done(response => {
                $(".page-content").html('<p>' + response + '</p>')
            })
        },
        remove: () => {
            clearPageContent()
        }
    },
    "sidebarPayment": {
        icon: "payment",
        label: "Payment",
        render: () => {
            const view = new ccPayment()
            view.render()
        },
        remove: () => {
            clearPageContent()
        }
    }
}

$.each(loaders, (key, loader) => {
    let $link = $(`<a id="${key}" class="mdl-navigation__link" href="#${key}"><i class="mdl-color-text--blue-grey-400 material-icons" role="presentation">${loader.icon}</i>${loader.label}</a>`)
    $link.bind('click', e => {
        if (window._currentLoaderId) loaders[window._currentLoaderId].remove()
        window._currentLoaderId = key
        loader.render()
    })
    $(".mdl-navigation").append($link)
})
