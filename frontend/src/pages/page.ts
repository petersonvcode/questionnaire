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
