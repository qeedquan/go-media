#include "stbte.h"
#include "_cgo_export.h"

int STB_TEXTEDIT_STRINGLEN(STB_TEXTEDIT_STRING *str) {
	return stringlen(str);
}

int STB_TEXTEDIT_GETCHAR(STB_TEXTEDIT_STRING *str, int idx) {
	return getchar(str, idx);
}

float STB_TEXTEDIT_GETWIDTH(STB_TEXTEDIT_STRING *str, int line_start_idx, int char_idx) {
	return getwidth(str, line_start_idx, char_idx);
}

void STB_TEXTEDIT_LAYOUTROW(StbTexteditRow *r, STB_TEXTEDIT_STRING *str, int line_start_idx) {
	return layoutrow(r, str, line_start_idx);
}

int STB_TEXTEDIT_MOVEWORDRIGHT(STB_TEXTEDIT_STRING *str, int idx) {
	return movewordright(str, idx);
}

int STB_TEXTEDIT_MOVEWORDLEFT(STB_TEXTEDIT_STRING *str, int idx) {
	return movewordleft(str, idx);
}

void STB_TEXTEDIT_DELETECHARS(STB_TEXTEDIT_STRING *str, int pos, int n) {
	deletechars(str, pos, n);
}

int STB_TEXTEDIT_INSERTCHARS(STB_TEXTEDIT_STRING *str, int pos, STB_TEXTEDIT_CHARTYPE *new_text, int new_text_len) {
	if (insertchars(str, pos, new_text, new_text_len))
		return 1;
	return 0;
}

int STB_TEXTEDIT_KEYTOTEXT(int key) {
	return key;
}
