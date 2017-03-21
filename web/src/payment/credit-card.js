import './user.scss'
import './credit-card.scss'
import userTemplate from './user.html'
import ccTemplate from './credit-card.html'
import {applyTemplate} from 'common/utils'

export default class ccPayment {
    constructor () {
        this.userFormId = 'stripe-user'
        this.ccFormId = 'stripe-cc'
    }
    render () {
        applyTemplate(".page-content", userTemplate({id: this.userFormId}))
        let $userForm = $(`#${this.userFormId}`)

        $userForm.find('#next').on('click', event => {
            event.preventDefault()
            applyTemplate(`#${this.userFormId} .payment`, ccTemplate({id: this.ccFormId}), true)
            $userForm.find('.payment').css('visibility', 'visible')
            $userForm.find('.submit').prop('disabled', false)
        })

        $userForm.submit(event => {
            let $ccForm = $(`#${this.ccFormId}`)

            // Disable the submit button to prevent repeated clicks:
            $userForm.find('.submit').prop('disabled', true)

            // Request a token from Stripe:
            Stripe.card.createToken($ccForm, this.stripeResponseHandler.bind(this))

            // Prevent the form from being submitted:
            return false
        })
    }
    stripeResponseHandler (status, response) {
        // Grab the form:
        let $ccForm = $(`#${this.ccFormId}`)
        let $userForm = $(`#${this.userFormId}`)

        if (response.error) { // Problem!

            // Show the errors on the form:
            $ccForm.find('.payment-errors').text(response.error.message)
            $userForm.find('.submit').prop('disabled', false) // Re-enable submission

        } else { // Token was created!
            let userData = this.prepareUserDataJson($userForm.serializeArray())
            userData['payment_service'] = "Stripe"
            userData['payment_token'] = response.id
            
            $.ajax({
                type: "POST",
                url: '/users',
                data: JSON.stringify(userData),
                contentType: 'application/json'
            }).done(response => {
                console.log(response)
            })
        }
    }
    prepareUserDataJson (userData) {
        const userJSON = {}
        $.each(userData, (idx, data) => {
            userJSON[data.name] = data.value
        })
        return userJSON
    }
}