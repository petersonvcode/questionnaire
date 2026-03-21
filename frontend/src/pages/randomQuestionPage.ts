import { Page, showElement, hideElement } from "./page";
import { getRandomQuestion, Question} from '../questions'

export class RandomQuestionPage implements Page {
  private readonly pageElement: HTMLElement;
  private readonly questionTextElement: HTMLElement;
  private readonly optionsContainerElement: HTMLElement;
  private readonly confirmButtonElement: HTMLButtonElement;
  private readonly anotherQuestionButtonElement: HTMLButtonElement;

  private question: Question | null = null;
  private selectedOptionIndex: number | null = null;

  constructor() {
    const pageElement = document.getElementById('random-question-page');
    if (!pageElement)
      throw new Error('Page element not found');
    this.pageElement = pageElement;

    const questionTextElement = document.getElementById('question-text');
    if (!questionTextElement)
      throw new Error('Question text element not found');
    this.questionTextElement = questionTextElement;

    const optionsContainerElement = document.querySelector('#random-question-page .options-container') as HTMLElement;
    if (!optionsContainerElement)
      throw new Error('Options container element not found');
    this.optionsContainerElement = optionsContainerElement;

    const confirmButtonElement = document.getElementById('confirm-rnd-question-btn');
    if (!confirmButtonElement || !(confirmButtonElement instanceof HTMLButtonElement))
      throw new Error('Confirm button element not found');
    this.confirmButtonElement = confirmButtonElement;
    this.confirmButtonElement.addEventListener('click', this.confirmQuestion);

    const anotherQuestionButtonElement = document.getElementById('another-rnd-question-btn');
    if (!anotherQuestionButtonElement || !(anotherQuestionButtonElement instanceof HTMLButtonElement))
      throw new Error('Another question button element not found');
    this.anotherQuestionButtonElement = anotherQuestionButtonElement;
    this.anotherQuestionButtonElement.addEventListener('click', () => {
      this.hideAnotherQuestionButton()
      this.loadQuestion()
    });
  }

  show = () => {
    showElement(this.pageElement)
    this.loadQuestion()
  }

  private loadQuestion = async () => {
    this.displayConfirmButton()
    this.questionTextElement.textContent = 'Carregando...';
    this.optionsContainerElement.innerHTML = '';
    this.question = await getRandomQuestion()
    this.questionTextElement.textContent = this.question.text;

    const buttons = this.question.options.map((option, i) => {
      const input = document.createElement('input');
      input.type = 'radio';
      input.name = 'q-op';
      input.value = option.id.toString();
      input.id = `q-op-${i}`;
      input.classList.add('q-op-input');
      input.addEventListener('change', () => {
        this.selectedOptionIndex = i;
        this.enableConfirmButton();
      });

      const label = document.createElement('label');
      label.htmlFor = `q-op-${i}`;
      label.appendChild(input);
      label.textContent = option.text;

      const container = document.createElement('div');
      container.className = 'q-op';
      container.appendChild(label);
      container.appendChild(input);
      return container;
    });
    for (const button of buttons)
      this.optionsContainerElement.appendChild(button);
  }

  private enableConfirmButton = () => {
    this.confirmButtonElement.disabled = false;
    this.confirmButtonElement.style.display = 'block';
  }

  private hideConfirmButton = () => {
    this.confirmButtonElement.style.display = 'none';
  }

  private displayConfirmButton = () => {
    this.confirmButtonElement.disabled = true;
    this.confirmButtonElement.style.display = 'block';
  }

  private confirmQuestion = () => {
    const allOptions = document.querySelectorAll('.q-op-input') as NodeListOf<HTMLInputElement>;
    for (const option of allOptions) {
      if (!("value" in option)) {
        console.warn('Option is not a radio input');
        continue;
      }
      const qOption = this.question?.options.find(o => o.id === parseInt(option.value as string));
      option.parentElement?.classList.add(qOption?.correct ? 'q-op-correct' : 'q-op-incorrect');
      option.disabled = true
    }

    this.hideConfirmButton()
    this.showAnotherQuestionButton()
  }

  private showAnotherQuestionButton = () => showElement(this.anotherQuestionButtonElement, 'block')

  private hideAnotherQuestionButton = () => hideElement(this.anotherQuestionButtonElement)

  hide = () => {
    hideElement(this.pageElement)
    this.question = null;
    this.questionTextElement.textContent = 'Carregando...';
    this.optionsContainerElement.innerHTML = '';
  }
}