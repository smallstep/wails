// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Call} from '@wailsio/runtime';
import {Person} from './models';

// Greet does XYZ
export async function Greet(name: string) : Promise<string> {
	return Call.ByID(1411160069, name);
}

// NewPerson creates a new person
export async function NewPerson(name: string) : Promise<Person> {
	return Call.ByID(1661412647, name);
}

