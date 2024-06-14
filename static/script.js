// static/script.js

const INTERVAL_DURATION = 1000;

class Clock {
  constructor() {
    this.clockElement = document.getElementById('app');
    this.updateTime();
    this.startClock();
  }

  startClock() {
    setInterval(
      () => {
        this.updateTime();
      }, 
      INTERVAL_DURATION,
    );
  }

  updateTime() {
    const now = new Date();
    const time = this.formatTime(now);
    const date = this.formatDate(now);
  
    if (this.clockElement) {
      this.clockElement.innerHTML = `
        <p>
          ${date} ${time}
        </p>
      `;
    }
  }
  
  formatTime(date) {
    const hours = this.padWithZero(date.getHours());
    const minutes = this.padWithZero(date.getMinutes());
    const seconds = this.padWithZero(date.getSeconds());
    return `${hours}:${minutes}:${seconds}`;
  }
  
  formatDate(date) {
    const year = date.getFullYear();
    const month = this.padWithZero(date.getMonth() + 1);
    const day = this.padWithZero(date.getDate());
    return `${year}/${month}/${day}`;
  }
  
  padWithZero(value) {
    return value.toString().padStart(2, '0');
  }

  getTime() {
    if (this.clockElement) {
      const timeText = this.clockElement.querySelector('p:nth-child(1)').innerText;
      return timeText.split(': ')[1];
    }
    return '';
  }

  getDate() {
    if (this.clockElement) {
      const dateText = this.clockElement.querySelector('p:nth-child(2)').innerText;
      return dateText.split(': ')[1];
    }
    return '';
  }
}

export default Clock;

window.addEventListener(
  'DOMContentLoaded', 
  function () {
    const clock = new Clock();
    clock.updateTime();
  },
);
