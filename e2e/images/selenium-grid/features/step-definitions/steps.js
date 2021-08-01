const { Given, When, Then } = require('@cucumber/cucumber');

Given(/^I am on the (\w+) page$/, async (page) => {
    await browser.url(`https://the-internet.herokuapp.com/${page}`);
});

When(/^I login with (\w+) and (.+)$/, async (username, password) => {
    await (await $('#username')).setValue(username);
    await (await $('#password')).setValue(password);
    await (await $('button[type="submit"]')).click();
});

Then(/^I should see a flash message saying (.*)$/, async (message) => {
    await expect($('#flash')).toBeExisting();
    await expect($('#flash')).toHaveTextContaining(message);
});

