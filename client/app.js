// TODO: scrape linkedIn profile https: //www.linkedin.com/pulse/how-easy-scraping-data-from-linkedin-profiles-david-craven/

class App {
  constructor () {
    this.inTraining = false;
    this.answerNode = document.querySelector('.answer');
    this.questionNode = document.querySelector('.question__text');
    this.tagNode = document.querySelector('.tags');
    this.selectNode = document.querySelector('.tags select');
    document.addEventListener('DOMContentLoaded', this.populateValues, false);
  }

  clearHighlight() {
    let highlightNodes = document.querySelectorAll('.highlight');
    highlightNodes.forEach(node => node.classList.remove('highlight'));
  }

  easterEgg() {
    window.location.href = 'https://www.intechnic.com/blog/best-website-easter-eggs-hidden-on-the-internet/';
  }

  // toggleButtonDisable() {
  //   let button = document.querySelector('.button');
  //   if(this.inTraining && button.disabled) {
  //     button.disabled = false;
  //   }
  // }

  findAnswer() {
    const dictKeys = Object.keys(answerDict)
    const randomIndex = Math.floor(Math.random() * dictKeys.length);
    return answerDict[dictKeys[randomIndex]];
  }

  onBlur() {
    if(!this.inTraining) {
      this.questionNode.value = '';
    }
  }

  onSubmit() {
    this.clearHighlight();
    if(!this.inTraining) {
      document.querySelectorAll('li').forEach(el => el.classList.remove('highlight'));
      const { message, el } = this.findAnswer(this.questionNode.value);
  
      if (el) {
        const elNode = document.getElementById(el);
        elNode.classList.add('highlight')
      }
      if (this.answerNode.innerHTML) {
        this.answerNode.innerHTML = '';
      }
      this.typeText(message);
    }
    return false;
  }

  openMail() {
    window.open('mailto:housker@gmail.com');
  }

  populateValues() {
    this.tooltipNode = document.querySelector('.tooltip__text');
    this.tooltipNode.innerHTML = INFO;
  }

  async populateTags() {
    const headers = new Headers({
      "Content-Type": "application/json",
      Accept: "application/json"
    });

    const tagsPromise = await fetch(`/tags`, headers);
    const tags = await tagsPromise.json();
    console.log('tags json: ', tags)

    tags.forEach(tag => {
      let el = document.createElement('option');
      el.innerHTML = tag;
      this.selectNode.appendChild(el);
    });
  }

  toggleDisplay(node) {
    if (node.classList.contains('nodisplay')) {
      node.classList.remove('nodisplay');
    } else {
      node.classList.add('nodisplay');
    }
  }

  async toggleTraining() {
    this.answerNode.innerHTML = '';
    this.questionNode.value = '';
    this.clearHighlight();
    this.inTraining = !this.inTraining;
    if(this.inTraining) {
      if(!this.selectNode.innerHTML.trim()) {
        await this.populateTags();
      }
    }
    let trainingNodes = document.querySelectorAll('.training');
    trainingNodes.forEach(node => {
      this.toggleDisplay(node);
    })
  }

  typeText(textToType) {
    const textArray = [...textToType];
    textArray.forEach((char, i) => {
      (function(index, context) {
        setTimeout(function() {
          context.answerNode.innerHTML += char;
          }, i * 70);
      })(i, this);
    });
  }

  updateCorpus() {
    const validMatches = this.validate();
    console.log('validMatches: ', validMatches, validMatches.length)
    if(validMatches.length < 0) {
      return;
    }
    console.log(this.selectNode.value);
    this.questionNode.value = '';
    // let floaterNode = document.querySelector('.floater');
    // floaterNode.classList.remove('nodisplay');
    // floaterNode.classList.add('floater__begin');
  }

  validate() {
    const regex = /\w+/gi;
    const matches = this.questionNode.value.matchAll(regex);
    return [...matches].map(match => match[0]);
  }
}

const app = new App();