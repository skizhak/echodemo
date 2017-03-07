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
        applyTemplate(".page-content", ccTemplate({id: this.ccFormId}), true)
        let $ccForm = $(`#${this.ccFormId}`)
        $ccForm.submit((event) => {
            // Disable the submit button to prevent repeated clicks:
            $ccForm.find('.submit').prop('disabled', true)

            // Request a token from Stripe:
            Stripe.card.createToken($ccForm, this.stripeResponseHandler.bind(this))

            // Prevent the form from being submitted:
            return false
        })
    }
    stripeResponseHandler (status, response) {
        // Grab the form:
        let $ccForm = $(`#${this.ccFormId}`)

        if (response.error) { // Problem!

            // Show the errors on the form:
            $ccForm.find('.payment-errors').text(response.error.message)
            $ccForm.find('.submit').prop('disabled', false) // Re-enable submission

        } else { // Token was created!

            // Get the token ID:
            let token = response.id

            // Insert the token ID into the form so it gets submitted to the server:
            $ccForm.append($('<input type="hidden" name="stripeToken">').val(token))

            // Submit the form:
            // $ccForm.get(0).submit()
        }
    }
}