import { Page, showElement, hideElement } from "./page";
import { Question } from "../questions";

export class AnsweringPage implements Page {
  private readonly pageElement: HTMLElement;
  private readonly questionNavigationElement: HTMLElement;
  private readonly questionTextElement: HTMLElement;
  private readonly questionOptionsContainerElement: HTMLElement;
  private readonly confirmAnsweringButtonElement: HTMLButtonElement;
  private readonly changeAnswerButtonElement: HTMLButtonElement;
  private readonly nextQuestionButtonElement: HTMLButtonElement;
  private readonly finishAnsweringButtonElement: HTMLButtonElement;
  private questionNavigationItems: HTMLDivElement[] = [];

  private questions: Question[] = [];
  // in questionAnswers the index is the same as the index in questions
  // and the value is the index of the selected option
  // if the value is -1, the question was not answered
  private questionAnswers: number[] = []
  private currentQuestionIndex: number = 0;

  constructor(
    onFinishAnswering: (
      questions: Question[],
      answersIndexes: number[]
    ) => void
  ) {
    const pageElement = document.getElementById('questions-answering-page');
    if (!pageElement)
      throw new Error('Page element not found');
    this.pageElement = pageElement;

    const questionNavigationElement = document.querySelector('#questions-answering-page .question-navigation') as HTMLElement;
    if (!questionNavigationElement)
      throw new Error('Question navigation element not found');
    this.questionNavigationElement = questionNavigationElement;

    const questionTextElement = document.getElementById('a-question-text');
    if (!questionTextElement)
      throw new Error('Question text element not found');
    this.questionTextElement = questionTextElement;

    const questionOptionsContainerElement = document.querySelector('#questions-answering-page .options-container') as HTMLElement;
    if (!questionOptionsContainerElement)
      throw new Error('Question options container element not found');
    this.questionOptionsContainerElement = questionOptionsContainerElement;
    
    const confirmAnsweringButtonElement = document.getElementById('confirm-answering-btn');
    if (!confirmAnsweringButtonElement || !(confirmAnsweringButtonElement instanceof HTMLButtonElement))
      throw new Error('Confirm answering button element not found');
    this.confirmAnsweringButtonElement = confirmAnsweringButtonElement;
    this.confirmAnsweringButtonElement.addEventListener('click', this.confirmQuestion);
    
    const nextQuestionButtonElement = document.getElementById('next-question-btn');
    if (!nextQuestionButtonElement || !(nextQuestionButtonElement instanceof HTMLButtonElement))
      throw new Error('Next question button element not found');
    this.nextQuestionButtonElement = nextQuestionButtonElement;
    this.nextQuestionButtonElement.addEventListener('click', this.goToNextQuestion);

    const changeAnswerButtonElement = document.getElementById('change-answer-btn');
    if (!changeAnswerButtonElement || !(changeAnswerButtonElement instanceof HTMLButtonElement))
      throw new Error('Change answer button element not found');
    this.changeAnswerButtonElement = changeAnswerButtonElement;
    this.changeAnswerButtonElement.addEventListener('click', this.changeAnswer);
    
    const finishAnsweringButtonElement = document.getElementById('finish-answering-btn');
    if (!finishAnsweringButtonElement || !(finishAnsweringButtonElement instanceof HTMLButtonElement))
      throw new Error('Finish answering button element not found');
    this.finishAnsweringButtonElement = finishAnsweringButtonElement;
    this.finishAnsweringButtonElement.addEventListener('click', () => onFinishAnswering(this.questions, this.questionAnswers));
  }

  private changeAnswer = () => {
    this.disableFinishAnsweringButton()
    this.questionAnswers[this.currentQuestionIndex] = -1;
    this.changeNavigationButtonColor(this.currentQuestionIndex, 'unanswered');
    this.loadQuestion(this.currentQuestionIndex);
  }

  private confirmQuestion = () => {
    const allOptions = this.questionOptionsContainerElement.querySelectorAll('input') as NodeListOf<HTMLInputElement>;
    const answerIndex = Array.from(allOptions).findIndex(o => o.checked);
    this.questionAnswers[this.currentQuestionIndex] = answerIndex;
    this.changeNavigationButtonColor(this.currentQuestionIndex, 'answered');

    const allAnswered = this.questionAnswers.every(a => a !== -1);
    if (allAnswered)
      this.enableFinishAnsweringButton()
    this.goToNextQuestion();
  }

  private goToNextQuestion = () => {
    const nextUnansweredIndex = this.questions.findIndex((_, index) => index > this.currentQuestionIndex && this.questionAnswers[index] === -1);
    const someUnanswered = this.questions.findIndex((_, index) => this.questionAnswers[index] === -1);
    const nextQuestionIndex = this.currentQuestionIndex + 1 >= this.questions.length  ? 0 : this.currentQuestionIndex + 1;

    const next = nextUnansweredIndex !== -1 ? nextUnansweredIndex 
      : someUnanswered !== -1 ? someUnanswered 
        : nextQuestionIndex;
    this.loadQuestion(next);
  }

  loadQuestions = (questions: Question[]) => {
    this.questions = questions;
    this.questionAnswers = new Array(questions.length).fill(-1);
    this.drawQuestionNavigation();
    this.loadQuestion(this.currentQuestionIndex);
  }

  private drawQuestionNavigation = () => {
    this.questionNavigationElement.innerHTML = '';
    const questionNavigationItems: HTMLDivElement[] = [];
    for (let i = 0; i < this.questions.length; i++) {
      const div = document.createElement('div');
      div.classList.add('question-navigation-item');
      div.textContent = `${i + 1}`;
      div.addEventListener('click', () => this.loadQuestion(i));
      this.questionNavigationElement.appendChild(div);
      questionNavigationItems.push(div);
    }
    this.questionNavigationItems = questionNavigationItems;
  }

  private loadQuestion = (index: number) => {
    this.currentQuestionIndex = index;
    this.changeNavigationButtonColor(this.currentQuestionIndex, 'selected');
    const question = this.questions[index];
    this.questionTextElement.textContent = question.text;
    this.questionOptionsContainerElement.innerHTML = '';

    const isAnswered = this.questionAnswers[index] !== -1;
    const answerIndex = this.questionAnswers[index];
    let optionIndex = 0
    for (const option of question.options) {
      const div = document.createElement('div');
      div.classList.add('question-option');

      const input = document.createElement('input');
      input.type = 'radio';
      input.name = 'q-op';
      input.value = option.id.toString();
      input.id = `q-op-${option.id}`;
      input.classList.add('q-op-input');
      input.addEventListener('change', this.enableConfirmButton);
      if (isAnswered) {
        input.checked = answerIndex === optionIndex;
        input.disabled = true;
      }

      const label = document.createElement('label');
      label.htmlFor = `q-op-${option.id}`;
      label.appendChild(input);
      label.textContent = option.text;

      div.appendChild(label);
      div.appendChild(input);

      this.questionOptionsContainerElement.appendChild(div);
      optionIndex++
    }
    if (isAnswered) {
      this.hideConfirmButton()
      this.displayChangeAnswerButton()
    }
    else {
      this.hideChangeAnswerButton()
      this.displayConfirmButton()
      this.disableConfirmButton()
    }
  }

  private enableConfirmButton = () => {
    this.confirmAnsweringButtonElement.disabled = false
  }

  private disableConfirmButton = () => {
    this.confirmAnsweringButtonElement.disabled = true;
  }

  private displayConfirmButton = () => {
    this.confirmAnsweringButtonElement.style.display = 'block';
  }

  private hideConfirmButton = () => {
    this.confirmAnsweringButtonElement.style.display = 'none';
  }

  private displayChangeAnswerButton = () => {
    this.changeAnswerButtonElement.style.display = 'block';
  }

  private hideChangeAnswerButton = () => {
    this.changeAnswerButtonElement.style.display = 'none';
  }

  private enableFinishAnsweringButton = () => {
    this.finishAnsweringButtonElement.disabled = false;
    this.finishAnsweringButtonElement.style.display = 'block';
  }

  private disableFinishAnsweringButton = () => {
    this.finishAnsweringButtonElement.disabled = true;
    this.finishAnsweringButtonElement.style.display = 'none';
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
  
  show = () => showElement(this.pageElement)

  hide = () => {
    hideElement(this.pageElement)
    this.disableFinishAnsweringButton()
    this.disableConfirmButton()
    this.questions = [];
    this.currentQuestionIndex = 0;
    this.questionNavigationElement.textContent = 'Carregando...';
    this.questionTextElement.textContent = 'Carregando...';
    this.questionOptionsContainerElement.textContent = 'Carregando...';
  }
}