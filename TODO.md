# MVP Ecommerce Site to Launch

## TODO Backend
*It's important to remember that we are phase 1 and MVP so we are not focusing on features such as membership loyalty, email campaigns, etc.*. We are launching an MVP that will allow us to sell products online, email users of orders, etc.

- [*] Finish all gets for Products, Product variation and Inventory check for the checkout process.
- [*] Implement Stripe Confirmation for payment processing.
- [ ] Add Redis to cache Cart items and cart sessions.? Do we even need to use Redis? Can we use in memory cache? Maybe Phase 2
- [ ] Organize DB layer endpoint into different files to be more organized.
- [*] Clean up the dockerfiles to be two stage builds for better performance, security and size. (Huge improvements! Whooo!!)
- [ ] Make sure search is working properly with database. Figure out if I need to add a full text index, other performance improvements. Phase one, we may just need to do simple searches based on names, maybe tags and what not.
- [ ] Checkout Process:
    - [ ] Get appropriate names for everything, into metadata
    - [ ] Completion of checkout process:
       - [*] Get Product details (Product and Variation) from DB layer.
       - [*] Get Inventory check for the checkout process.
       - [*] If inventory is not available, notify the user and do not proceed to finish payments.
       - [++] Implement Strip confirmation for payment processing. Including Removing inventory from the database, sending email, updating user account with an order record (user can login to see orders progress, history and later maybe tracking?). Setup up local webhook testing.
       - [ ] Checkout-endpoint  is now going to be in the  c
- [ ] Process Email:
    - [ ] Research how to send emails with Go. Use SMTP or a third party service like Mailgun, SendGrid, Amazon SES.
    - [ ] Send email when order is placed
    - [ ] Send email when order is shipped
    - [ ] Send email when order is delivered
- [ ] User Account Features
    - [ ] User orders endpoints
- [-] Orders System
    - [ ] Implement a system to keep track of order progress
- [ ] Admin
- [ ] Envrinment Variables handling, better!
    -[*] Stripe
    -[ ] Hosts



## Frontend Improvements:
- [ ] Implement Product to Product variation view
- [ ] Implement Succesful purchase handling and display of order summary.
- [ ] User
    - [ ] Implement Order History and Tracking for User
    - [ ] Order Summary/Info Modal or page? 
        - [ ] Order summary information (Products, paid, address, etc)
        - [ ] Order status (Processing, shipped, delivered)
        - [ ] Shipping details (Shipping Tracking/Link)
- [ ] Admin features
