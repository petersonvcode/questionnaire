export type Question = {
  id: number;
  text: string;
  type: 'single_choice' | 'multiple_choice';
  tags: string[];
  options: {
    id: number;
    text: string;
    correct: boolean;
  }[];
}

export type Tag = {
  id: number;
  tag: string;
  count: number;
}

//@ts-ignore
const baseUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:9090';

export const getRandomQuestion = async (): Promise<Question> => {
  const url = `${baseUrl}/questions`;

  const response = await fetch(url);
  if (!response.ok)
    throw new Error('Failed to fetch random questions');
  const question = await response.json() as Question;
  return question;
}

export const getTags = async (): Promise<Tag[]> => {
  const url = `${baseUrl}/questions/tags`;
  const response = await fetch(url);
  if (!response.ok)
    throw new Error('Failed to fetch tags');
  const tags = await response.json() as Tag[];
  return tags;
}

export type OtherTagCount = {
  tagID: number;
  count: number;
}

export type CountQuestionsWithTags = {
  totalCount: number;
  otherTagsCount: OtherTagCount[];
}

export const countQuestionsWithTags = async (tags: Tag[]): Promise<CountQuestionsWithTags> => {
  const ids = tags.map(t => t.id);
  const idsStr = ids.join(',');
  const url = `${baseUrl}/questions/tags/count?ids=${idsStr}`;
  const response = await fetch(url);
  if (!response.ok)
    throw new Error('Failed to count questions with tags');
  return await response.json() as CountQuestionsWithTags;
}

export const getQuestionsWithTags = async (tagIDs: number[], count: number): Promise<Question[]> => {
  const idsStr = tagIDs.join(',');
  const url = `${baseUrl}/questions?tagIDs=${idsStr}&count=${count}`;
  const response = await fetch(url);
  if (!response.ok)
    throw new Error('Failed to get questions with tags');
  const questions = await response.json() as Question[];
  return questions;
}