export type Question = {
  id: string;
  text: string;
  type: 'single-choice';
  tags: string[];
  options: {
    id: string;
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
    input.value = option.id;
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

export const getRandomQuestion = async (): Promise<Question> => {
  // TODO change to actual API call
  const rndWait = Math.floor(Math.random() * 1000) + 500;
  await new Promise(resolve => setTimeout(resolve, rndWait));
  
  const questions: Question[] = [
    {
      id: '1',
      text: 'What is the capital of France?',
      type: 'single-choice',
      tags: ['geography', 'world'],
      options: [
        { id: '1', text: 'Paris', correct: true },
        { id: '2', text: 'London', correct: false },
        { id: '3', text: 'Berlin', correct: false },
      ],
    },
    {
      id: '2',
      text: 'Which planet is known as the Red Planet?',
      type: 'single-choice',
      tags: ['science', 'space'],
      options: [
        { id: '1', text: 'Venus', correct: false },
        { id: '2', text: 'Mars', correct: true },
        { id: '3', text: 'Jupiter', correct: false },
      ],
    },
    {
      id: '3',
      text: 'Who wrote "Romeo and Juliet"?',
      type: 'single-choice',
      tags: ['literature', 'authors'],
      options: [
        { id: '1', text: 'William Shakespeare', correct: true },
        { id: '2', text: 'Charles Dickens', correct: false },
        { id: '3', text: 'Jane Austen', correct: false },
      ],
    },
    {
      id: '4',
      text: 'What is the chemical symbol for water?',
      type: 'single-choice',
      tags: ['science', 'chemistry'],
      options: [
        { id: '1', text: 'O2', correct: false },
        { id: '2', text: 'CO2', correct: false },
        { id: '3', text: 'H2O', correct: true },
      ],
    },
    {
      id: '5',
      text: 'Which continent is Egypt located in?',
      type: 'single-choice',
      tags: ['geography', 'world'],
      options: [
        { id: '1', text: 'Asia', correct: false },
        { id: '2', text: 'Africa', correct: true },
        { id: '3', text: 'Europe', correct: false },
      ],
    },
    {
      id: '6',
      text: 'What is the largest mammal in the world?',
      type: 'single-choice',
      tags: ['biology', 'animals'],
      options: [
        { id: '1', text: 'Elephant', correct: false },
        { id: '2', text: 'Blue Whale', correct: true },
        { id: '3', text: 'Giraffe', correct: false },
      ],
    },
    {
      id: '7',
      text: 'In what year did the World War II end?',
      type: 'single-choice',
      tags: ['history', '20th century'],
      options: [
        { id: '1', text: '1945', correct: true },
        { id: '2', text: '1939', correct: false },
        { id: '3', text: '1918', correct: false },
      ],
    },
    {
      id: '8',
      text: 'Which element has the atomic number 1?',
      type: 'single-choice',
      tags: ['science', 'chemistry'],
      options: [
        { id: '1', text: 'Helium', correct: false },
        { id: '2', text: 'Hydrogen', correct: true },
        { id: '3', text: 'Carbon', correct: false },
      ],
    }
  ];
  
  const randomIndex = Math.floor(Math.random() * questions.length);
  return questions[randomIndex];
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
    const qOption = currentQuestion?.options.find(o => o.id === option.value);

    option.parentElement?.classList.add(qOption?.correct ? 'q-op-correct' : 'q-op-incorrect');
    disableConfirmButton()
  }
}