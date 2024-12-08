package runes

import "unicode"

type UnderscoreOptions struct {
	PreserveCase bool
	Screaming    bool
}

func Underscore(runes []rune, options *UnderscoreOptions) []rune {
	if len(runes) == 0 {
		return runes
	}

	sb := make([]rune, 0)
	last := rune(0)
	if options == nil {
		options = &UnderscoreOptions{}
	}

	for _, r := range runes {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				if unicode.IsLetter(last) && unicode.IsLower(last) {
					sb = append(sb, '_')
					if options.PreserveCase || options.Screaming {
						sb = append(sb, r)
						last = r
						continue
					}

					sb = append(sb, unicode.ToLower(r))
					last = r
					continue
				}

				if options.PreserveCase || options.Screaming {
					sb = append(sb, r)
					last = r
					continue
				}

				sb = append(sb, unicode.ToLower(r))
				last = r
				continue
			}

			if options.Screaming {
				sb = append(sb, unicode.ToUpper(r))
			} else if options.PreserveCase {
				sb = append(sb, r)
			} else {
				sb = append(sb, unicode.ToLower(r))
			}

			last = r
			continue
		}

		if unicode.IsNumber(r) {
			sb = append(sb, r)
			last = r
			continue
		}

		if r == '_' || r == '-' || unicode.IsSpace(r) {
			if len(sb) == 0 {
				continue
			}

			if last == '_' {
				continue
			}

			last = '_'
			sb = append(sb, last)
			continue
		}

	}

	if len(sb) > 0 && sb[len(sb)-1] == '_' {
		sb = sb[:len(sb)-1]
	}

	return sb
}
