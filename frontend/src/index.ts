import { confirmQuestion, desiredQuestionsCount, getRandomQuestion, getTags, loadQuestion, loadTags, setDesiredQuestionsCount, unloadQuestion, unloadTags } from "./questions";
import './index.css';

document.addEventListener('DOMContentLoaded', () => {
  // Adding event listeners to main page buttons
  document.getElementById('main-random-btn')
    ?.addEventListener('click', renderRandomQuestionPage);

  const btns = document.getElementsByClassName('go-to-main-btn')
  for (const btn of btns)
    btn.addEventListener('click', renderMainPage);

  document.getElementById('confirm-question-btn')
    ?.addEventListener('click', confirmQuestion);

  document.getElementById('themes-btn')
    ?.addEventListener('click', renderSelectThemesPage);

  const cntBtns = document.getElementsByClassName('q-count-btn');
  for (const btn of cntBtns)
    btn.addEventListener('click', () => {
      const newCount = parseInt(btn.textContent || '0');
      if (isNaN(newCount)) {
        console.warn(`Invalid desired questions count: ${btn.textContent}`);
        return;
      }
      setDesiredQuestionsCount(newCount);
    });
})

const renderRandomQuestionPage = () => {
  hideAllPagesContent();
  showAnswersContainer();
  getRandomQuestion().then(q => loadQuestion(q, true)).catch(console.error);
}

const renderMainPage = () => {
  hideAllPagesContent();

  unloadQuestion();
  unloadTags();
  
  showMainButtonContainer();
}

const renderSelectThemesPage = () => {
  hideAllPagesContent();
  showSelectThemesContainer();
  getTags().then(loadTags).catch(console.error);
}

const hideAllPagesContent = () => {
  const elements = [
    document.getElementById('main-button-container'),
    document.getElementById('answers-container'),
    document.getElementById('select-themes-container'),
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
const showSelectThemesContainer = () => showElement(document.getElementById('select-themes-container'), 'flex');