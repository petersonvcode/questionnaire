export type Question = {
  id: number;
  text: string;
  type: 'single-choice';
  tags: string[];
  options: {
    id: number;
    text: string;
    correct: boolean;
  }[];
}

export let isRandomQuestion = false;
let currentQuestion: Question | null = null;

export const getCurrentQuestion = (): Question | null => currentQuestion;

export const loadQuestion = (question: Question, isRandom: boolean = false) => {
  isRandomQuestion = isRandom;
  currentQuestion = question;
  // Carregar valores no DOM
  const textElement = document.getElementById('question-text');
  if (!textElement)
    throw new Error('Question text element not found');
  textElement.textContent = question.text;

  const optionsElement = document.getElementById('options-container');
  if (!optionsElement)
    throw new Error('Options container element not found');

  // Creating options buttons
  const buttons = question.options.map((option, i) => {
    const input = document.createElement('input');
    input.type = 'radio';
    input.name = 'q-op';
    input.value = option.id.toString();
    input.id = `q-op-${i}`;
    input.classList.add('q-op-input');
    input.addEventListener('change', enableConfirmButton);

    const label = document.createElement('label');
    label.htmlFor = `q-op-${i}`;
    label.appendChild(input);
    label.textContent = option.text;

    const container = document.createElement('div');
    container.appendChild(label);
    container.appendChild(input);
    return container;
  })
  for (const button of buttons)
    optionsElement.appendChild(button);
}


export const unloadQuestion = () => {
  currentQuestion = null;
  isRandomQuestion = false;
  // Limpar valores no DOM
  const textElement = document.getElementById('question-text');
  if (textElement)
    textElement.textContent = 'Carregando...';
  const optionsElement = document.getElementById('options-container');
  if (optionsElement)
    optionsElement.innerHTML = '';

  disableConfirmButton()
}

//@ts-ignore
const baseUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080';

export const getRandomQuestion = async (): Promise<Question> => {
  const url = `${baseUrl}/questions`;

  const response = await fetch(url);
  if (!response.ok)
    throw new Error('Failed to fetch random questions');
  const question = await response.json() as Question;
  return question;
}

const enableConfirmButton = () => {
  const confirmButton = document.getElementById('confirm-question-btn');
  if (confirmButton && confirmButton instanceof HTMLButtonElement)
    confirmButton.disabled = false;
}

const disableConfirmButton = () => {
  const confirmButton = document.getElementById('confirm-question-btn');
  if (confirmButton && confirmButton instanceof HTMLButtonElement)
    confirmButton.disabled = true;
}

const removeConfirmButton = () => {
  const confirmButton = document.getElementById('confirm-question-btn');
  if (confirmButton && confirmButton instanceof HTMLButtonElement)
    confirmButton.remove();
}

const addConfirmButton = () => {
  const confirmButton = document.getElementById('confirm-question-btn');
  if (confirmButton && confirmButton instanceof HTMLButtonElement) {
    console.log('Confirm button already exists');
    return;
  }

  const button = document.createElement('button');
  button.id = 'confirm-question-btn';
  button.textContent = 'Confirmar';
  button.addEventListener('click', confirmQuestion);
  document.getElementById('buttons-container')?.appendChild(button);
}

const addNextQuestionButton = () => {
  const nextButton = document.getElementById('next-question-btn');
  if (nextButton && nextButton instanceof HTMLButtonElement) {
    console.log('Next button already exists');
    return;
  }

  const button = document.createElement('button');
  button.id = 'next-question-btn';
  button.textContent = 'Próxima questão';
  button.addEventListener('click', () => {
    unloadQuestion()
    removeNextQuestionButton()
    addConfirmButton()
    getRandomQuestion().then(q => loadQuestion(q, true)).catch(console.error);
  });
  document.getElementById('buttons-container')?.appendChild(button);
}

const removeNextQuestionButton = () => {
  const nextButton = document.getElementById('next-question-btn');
  if (nextButton && nextButton instanceof HTMLButtonElement)
    nextButton.remove();
}

export const confirmQuestion = () => {
  if (!isRandomQuestion)
    throw new Error('Not implemented');

  // Since this is a single random question, we can display
  // the correct answer in the options container itself
  const selectedOption = document.querySelector('.q-op-input:checked');
  if (!selectedOption || !("value" in selectedOption))
    throw new Error('No option selected');

  const allOptions = document.querySelectorAll('.q-op-input');
  for (const option of allOptions) {
    if (!("value" in option)) {
      console.warn('Option is not a radio input');
      continue;
    }
    const qOption = currentQuestion?.options.find(o => o.id === parseInt(option.value as string));

    option.parentElement?.classList.add(qOption?.correct ? 'q-op-correct' : 'q-op-incorrect');
    removeConfirmButton()
    addNextQuestionButton()
  }
}