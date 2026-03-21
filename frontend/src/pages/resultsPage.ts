import { Question } from "../questions";
import { hideElement, Page, showElement } from "./page";

export class ResultsPage implements Page {
  private readonly pageElement: HTMLElement;
  private readonly questionNavigationElement: HTMLElement;
  private questionNavigationItems: HTMLDivElement[] = [];
  private readonly questionTextElement: HTMLElement;
  private readonly questionOptionsContainerElement: HTMLElement;
  private readonly summaryElement: HTMLElement;
  private readonly nextButtonElement: HTMLButtonElement;

  private questionIndex = 0
  private questions: Question[] = [];
  private answersIndexes: number[] = [];

  constructor() {
    const pageElement = document.getElementById('results-page');
    if (!pageElement)
      throw new Error('Page element not found');

    const questionNavigationElement = document.querySelector('#results-page .question-navigation') as HTMLElement;
    if (!questionNavigationElement)
      throw new Error('Question navigation element not found');
    this.questionNavigationElement = questionNavigationElement;

    const questionTextElement = document.getElementById('r-question-text');
    if (!questionTextElement)
      throw new Error('Question text element not found');
    this.questionTextElement = questionTextElement;

    const questionOptionsContainerElement = document.querySelector('#results-page .options-container') as HTMLElement;
    if (!questionOptionsContainerElement)
      throw new Error('Question options container element not found');
    this.questionOptionsContainerElement = questionOptionsContainerElement;

    const nextButtonElement = document.getElementById('r-next-btn');
    if (!nextButtonElement || !(nextButtonElement instanceof HTMLButtonElement))
      throw new Error('Next button element not found');
    this.nextButtonElement = nextButtonElement;
    this.nextButtonElement.addEventListener('click', this.loadNextQuestion);

    const summaryElement = document.getElementById('r-summary');
    if (!summaryElement)
      throw new Error('Summary element not found');
    this.summaryElement = summaryElement;

    this.pageElement = pageElement;
  }

  show = () => showElement(this.pageElement)
  hide = () => {
    hideElement(this.pageElement)
    this.questions = [];
    this.answersIndexes = [];
    this.questionIndex = 0;
    this.questionNavigationItems = [];
  }

  private loadNextQuestion = () => {
    this.questionIndex++;
    if (this.questionIndex >= this.questions.length)
      this.questionIndex = 0;
    this.loadQuestion(this.questionIndex);
  }

  loadQuestions = (questions: Question[], answersIndexes: number[]) => {
    this.questions = questions;
    this.answersIndexes = answersIndexes;
    this.drawQuestionNavigation()
    this.loadQuestion(this.questionIndex);

    const correctAnswers = this.questions.reduce((acc, question, index) => {
      const answerIndex = this.answersIndexes[index];
      const isCorrect = question.options[answerIndex].correct;
      return acc + (isCorrect ? 1 : 0);
    }, 0);
    this.summaryElement.textContent = `${correctAnswers}/${this.questions.length} (${Math.round((correctAnswers / this.questions.length) * 100)}%)`;
  }

  private drawQuestionNavigation = () => {
    this.questionNavigationElement.innerHTML = '';
    for (let i = 0; i < this.questions.length; i++) {
      const answerIndex = this.answersIndexes[i];
      const isCorrect = this.questions[i].options[answerIndex].correct;
      const div = document.createElement('div');
      div.classList.add('question-navigation-item', isCorrect ? 'correct' : 'incorrect');
      div.textContent = `${i + 1}`;
      div.addEventListener('click', () => this.loadQuestion(i));
      this.questionNavigationItems.push(div);
      this.questionNavigationElement.appendChild(div);
    }
  }

  private loadQuestion = (index: number) => {
    this.questionIndex = index;
    this.changeNavigationButtonColor(this.questionIndex, 'selected');
    const question = this.questions[index];
    this.questionTextElement.textContent = question.text;

    this.questionOptionsContainerElement.innerHTML = '';
    let optionIndex = 0;
    for (const option of question.options) {
      const div = document.createElement('div');
      div.classList.add('question-option');
      div.classList.add(option.correct ? 'q-op-correct' : 'q-op-incorrect');
      
      const input = document.createElement('input');
      input.type = 'radio';
      input.name = 'q-op';
      input.value = option.id.toString();
      input.id = `q-op-${option.id}`;
      input.classList.add('q-op-input');
      input.disabled = true;
      const checked = this.answersIndexes[index] === optionIndex;
      input.checked = checked;

      const label = document.createElement('label');
      label.htmlFor = `q-op-${option.id}`;
      label.appendChild(input);
      label.textContent = option.text;

      div.appendChild(label);
      div.appendChild(input);

      this.questionOptionsContainerElement.appendChild(div);
      optionIndex++;
    }
  }

  private changeNavigationButtonColor = (
    questionIndex: number,
    type: 'selected' | 'answered' | 'unanswered' | 'correct' | 'incorrect'
  ) => {
    if (type === 'selected')
      for (const item of this.questionNavigationItems)
        item.classList.remove('selected')

    const item = this.questionNavigationItems[questionIndex];
    if (!item) {
      console.warn('Question navigation item not found');
      return;
    }

    if (type === 'unanswered')
      item.classList.remove('answered');
    item.classList.add(type);
  }
}