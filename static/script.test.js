// static/script.test.js

import Clock from './script.js';

test('Clock getTime method returns current time as string', () => {
    const clock = new Clock();
    clock.updateTime(); // Initialize the clock
    expect(typeof clock.getTime()).toBe('string');
});

test('Clock getDate method returns current date as string', () => {
    const clock = new Clock();
    clock.updateTime(); // Initialize the clock
    expect(typeof clock.getDate()).toBe('string');
});
