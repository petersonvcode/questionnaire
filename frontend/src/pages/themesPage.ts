import { hideElement, Page, showElement } from "./page";
import { countQuestionsWithTags, getTags, Tag } from "../questions";

export class ThemesQuestionPage implements Page {
  private readonly pageElement: HTMLElement;
  private readonly themesContainerElement: HTMLElement;
  private readonly themesListElement: HTMLElement;
  private readonly qAvailableElement: HTMLElement;
  private readonly qCountElements: HTMLButtonElement[];
  private readonly confirmTagsButtonElement: HTMLButtonElement;

  private selectedTags: Tag[] = [];
  private availableQuestionsCount: number | null = null;
  private desiredQuestionsCount: number | null = null;

  constructor(
    onConfirmTags: (tags: Tag[], count: number) => void
  ) {
    const pageElement = document.getElementById('themes-question-page');
    if (!pageElement)
      throw new Error('Page element not found');
    this.pageElement = pageElement;

    const themesContainerElement = document.getElementById('themes-container');
    if (!themesContainerElement)
      throw new Error('Themes container element not found');
    this.themesContainerElement = themesContainerElement;

    const themesListElement = document.getElementById('themes-list');
    if (!themesListElement)
      throw new Error('Themes list element not found');
    this.themesListElement = themesListElement;

    const qAvailableElement = document.getElementById('q-available');
    if (!qAvailableElement)
      throw new Error('Q available element not found');
    this.qAvailableElement = qAvailableElement;

    const confirmTagsButtonElement = document.getElementById('confirm-tags-btn');
    if (!confirmTagsButtonElement || !(confirmTagsButtonElement instanceof HTMLButtonElement))
      throw new Error('Confirm tags button element not found');
    this.confirmTagsButtonElement = confirmTagsButtonElement;
    this.confirmTagsButtonElement.addEventListener('click', () => onConfirmTags(this.selectedTags, this.desiredQuestionsCount ?? 0));

    const qCountBtns = document.getElementsByClassName('q-count-btn');
    this.qCountElements = Array.from(qCountBtns) as HTMLButtonElement[];
    for (const btn of this.qCountElements)
      btn.addEventListener('click', () => {
        this.changeDesiredQuestionsCount(parseInt(btn.textContent || '0'));
      });
  }

  show = () => {
    this.themesListElement.innerHTML = ''
    showElement(this.pageElement)
    showElement(this.themesContainerElement)
    getTags().then(tags => {
      for (const t of tags) {
        const div = document.createElement('div');
        div.classList.add('theme-item');
        div.textContent = t.tag;
        div.addEventListener('click', async () => {
          const isSelected = div.classList.contains('selected');

          if (isSelected) {
            div.classList.remove('selected')
            this.selectedTags = this.selectedTags.filter(s => s.id !== t.id);
          } else {
            div.classList.add('selected');
            this.selectedTags.push(t);
          }
          this.qAvailableElement.textContent = '...'
          const count = await countQuestionsWithTags(this.selectedTags)
          this.qAvailableElement.textContent = `Disponíveis: ${count}`;
          this.availableQuestionsCount = count;
          this.updateDesiredCountButtons()

          if (this.desiredQuestionsCount !== null && this.availableQuestionsCount !== null && this.desiredQuestionsCount > 0)
            this.enableConfirmButton()
          else
            this.disableConfirmButton()
        });
        this.themesListElement.appendChild(div);
      }
    })
  }

  hide = () => {
    hideElement(this.pageElement)
    hideElement(this.themesContainerElement)
    this.selectedTags = [];
    this.availableQuestionsCount = null;
    this.desiredQuestionsCount = null;
    this.qAvailableElement.textContent = 'Disponíveis: 0';
    for (const btn of this.qCountElements) {
      btn.disabled = true;
      btn.classList.remove('selected');
    }
    this.disableConfirmButton()
  }

  private enableConfirmButton = () => {
    this.confirmTagsButtonElement.disabled = false;
    this.confirmTagsButtonElement.style.display = 'block';
  }

  private disableConfirmButton = () => {
    this.confirmTagsButtonElement.disabled = true;
    this.confirmTagsButtonElement.style.display = 'block';
  }

  private changeDesiredQuestionsCount = (newCount: number) => {
    if (newCount > (this.availableQuestionsCount || 0)) {
      console.warn('Cannot set desired questions count to a value greater than the available questions count');
      return
    }
    this.desiredQuestionsCount = newCount;

    this.updateDesiredCountButtons()

    if (this.desiredQuestionsCount !== null && this.availableQuestionsCount !== null && this.desiredQuestionsCount > 0)
      this.enableConfirmButton()
    else
      this.disableConfirmButton()
  }

  private updateDesiredCountButtons = () => {
    for (const btn of this.qCountElements) {
      const btnValue = parseInt(btn.textContent || '0');
      if (btnValue > (this.availableQuestionsCount || 0)) {
        btn.disabled = true;
        if (btn.classList.contains('selected')) {
          btn.classList.remove('selected');
          this.desiredQuestionsCount = null;
        }
      } else
        btn.disabled = false;

      if (btnValue === (this.desiredQuestionsCount || 0))
        btn.classList.add('selected');
      else
        btn.classList.remove('selected');
    }
  }
}