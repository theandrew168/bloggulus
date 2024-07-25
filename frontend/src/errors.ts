// A structured error can either be:
// 1. General (just a message)
// 2. Specific (has a message that pertains to a specific field)

// Things that can happen:
// 1. All good, no errors
// 2. Validation errors (422)
//    1. Decode as an ErrorsResponse
//    2. Return to the form for fixing / resubmission
// 3. Non-validation errors
//    1. Try to parse as an ErrorsResponse
//    2. If it works, find the first message w/o a specific field
//    3. Otherwise, decode as plain text

export type StructuredError = {
	message: string;
	field?: string;
};

export type StructuredErrorsResponse = {
	errors: StructuredError[];
};

export function isStructuredErrorsResponse(value: any): value is StructuredErrorsResponse {
	return "errors" in value;
}

export function findGeneralError(errors: StructuredError[]): string | undefined {
	return errors.find((e) => !e.field)?.message;
}

export function findSpecificErrors(errors: StructuredError[]): Record<string, string> {
	const errorsByField = errors.reduce(
		(acc, err) => {
			if (err.field && !acc[err.field]) {
				acc[err.field] = err.message;
			}
			return acc;
		},
		{} as Record<string, string>,
	);
	return errorsByField;
}
