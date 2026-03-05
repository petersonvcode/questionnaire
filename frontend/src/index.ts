import { confirmQuestion, getRandomQuestion, loadQuestion, unloadQuestion } from "./questions";
import './index.css';

document.addEventListener('DOMContentLoaded', () => {
  // Adding event listeners to main page buttons
  document.getElementById('main-random-btn')
    ?.addEventListener('click', renderRandomQuestionPage);

  document.getElementById('go-to-main-btn')
    ?.addEventListener('click', renderMainPage);

  document.getElementById('confirm-question-btn')
    ?.addEventListener('click', confirmQuestion);
})

const renderRandomQuestionPage = () => {
  hideAllPagesContent();
  showAnswersContainer();
  getRandomQuestion().then(q => loadQuestion(q, true)).catch(console.error);
}

const renderMainPage = () => {
  hideAllPagesContent();
  unloadQuestion();
  showMainButtonContainer();
}

const hideAllPagesContent = () => {
  const elements = [
    document.getElementById('main-button-container'),
    document.getElementById('answers-container'),
  ]
  for (const element of elements)
    hideElement(element);
}

const hideElement = (element?: HTMLElement | null) => {
  if (element)
    element.style.display = 'none';
}

const showElement = (element?: HTMLElement | null, display: 'flex' | 'block' = 'flex') => {
  if (element)
    element.style.display = display;
}

const showAnswersContainer = () => showElement(document.getElementById('answers-container'), 'flex');
const showMainButtonContainer = () => showElement(document.getElementById('main-button-container'), 'flex');