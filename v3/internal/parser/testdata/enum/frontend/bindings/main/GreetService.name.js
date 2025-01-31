// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * @typedef {import('./models').Person} Person
 * @typedef {import('./models').Title} Title
 */

/**
 * Greet does XYZ
 * @function Greet
 * @param name {string}
 * @param title {Title}
 * @returns {Promise<string>}
 **/
export async function Greet(name, title) {
	return wails.CallByName("main.GreetService.Greet", ...Array.prototype.slice.call(arguments, 0));
}

/**
 * NewPerson creates a new person
 * @function NewPerson
 * @param name {string}
 * @returns {Promise<Person>}
 **/
export async function NewPerson(name) {
	return wails.CallByName("main.GreetService.NewPerson", ...Array.prototype.slice.call(arguments, 0));
}
