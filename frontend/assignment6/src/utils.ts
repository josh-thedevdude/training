import type { NewNote, Note } from "./main";

export function getNotes(): Note[] {
    const jsonString = localStorage.getItem("notes")
    if (!jsonString) return []
    try {
        return JSON.parse(jsonString);
    } catch (error) {
        console.error("Invalid JSON in localStorage for notes");
        return [];
    }
}

export function addNote(note: NewNote) {
    if (!note.title) {
        throw new Error("Title is required")
    }
    if (!note.description) {
        throw new Error("Description is required")
    }
    const newNote: Note = {
        ...note,
        createdAt: new Date().toDateString()
    }

    const notes = getNotes()
    notes.push(newNote);
    localStorage.setItem("notes", JSON.stringify(notes))
}

