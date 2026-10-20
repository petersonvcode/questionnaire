import { Question } from "../questions";

export type Page = {
  show(): void;
  hide(): void;
}

export const showElement = (element?: HTMLElement | null, display: 'flex' | 'block' = 'flex') => {
  if (element)
    element.style.display = display;
}

export const hideElement = (element?: HTMLElement | null) => {
  if (element)
    element.style.display = 'none';
}

export const flashAnswer = (correct: boolean) => {
  const el = document.getElementById('answer-flash');
  if (!el) return;
  el.textContent = correct ? '✓' : '✗';
  el.className = '';
  void el.offsetWidth;
  el.className = correct ? 'ok' : 'bad';
}

const optionInputClass = 'q-op-input';

export const displayQuestionAndOptions = (
  question: Question,
  questionTextElement: HTMLElement,
  questionOptionsContainerElement: HTMLElement,
  onOptionChange: (optionId: number) => void,
) => {
  questionTextElement.textContent = question.text;

  const optionsButtons = question.options.map((option, i) => {
    const input = document.createElement('input');
    input.type = question.type === 'single_choice' ? 'radio' : 'checkbox';
    input.name = 'q-op';
    input.value = option.id.toString();
    input.id = `q-op-${i}`;
    input.classList.add(optionInputClass);
    input.addEventListener('change', (e) => {
      onOptionChange(option.id);
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
  for (const button of optionsButtons)
    questionOptionsContainerElement.appendChild(button);
}

type ConfirmQuestionResult = {
  allCorrect: boolean;
  options: {
    id: number;
    userSelected: boolean;
    userGotItRight: boolean;
  }[];
}

export const confirmQuestion = (question: Question) => {
  const allOptions = document.getElementsByClassName(optionInputClass) as HTMLCollectionOf<HTMLInputElement>;
  const result: ConfirmQuestionResult = { allCorrect: true, options: [] };
  for (const option of allOptions) {
    if (!("value" in option)) {
      console.warn('Option input has no value');
      continue;
    }
    const optionId = parseInt(option.value as string);
    const qOption = question.options.find(o => o.id === optionId);
    if (!qOption) {
      console.warn('Option not found');
      continue;
    }
    if (!option.parentElement) {
      console.warn('Option parent element not found');
      continue;
    }

    option.parentElement.classList.add(qOption.correct ? 'q-op-correct' : 'q-op-incorrect');
    option.disabled = true

    const userSelected = option.checked;
    const userGotItRight = qOption.correct === userSelected;
    result.options.push({ id: optionId, userSelected, userGotItRight });
    if (!userGotItRight)
      result.allCorrect = false;
  }

  return result;
}