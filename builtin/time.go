package builtin

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"time"
)

// InitializeTimeFunctions инициализирует встроенные функции для работы с датой и временем
func InitializeTimeFunctions(globalScope *scope.ScopeStack) {
	// now - текущее время
	nowFunc := func(args []*value.Value) *value.Value {
		if len(args) != 0 {
			return value.NewString("Error: now() requires 0 arguments")
		}
		return value.NewTime(time.Now())
	}
	globalScope.Set("now", value.NewValue(nowFunc))

	// timeFromUnix - создание времени из Unix timestamp (секунды)
	timeFromUnixFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: timeFromUnix() requires 1 argument (unix timestamp)")
		}
		
		timestamp, ok := args[0].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[0].Any().(float64); ok {
				timestamp = int64(floatVal)
			} else {
				return value.NewString("Error: timeFromUnix() requires numeric argument")
			}
		}
		
		return value.NewTime(time.Unix(timestamp, 0))
	}
	globalScope.Set("timeFromUnix", value.NewValue(timeFromUnixFunc))

	// timeFromString - парсинг времени из строки
	timeFromStringFunc := func(args []*value.Value) *value.Value {
		if len(args) < 1 || len(args) > 2 {
			return value.NewString("Error: timeFromString() requires 1-2 arguments (timeString, [format])")
		}
		
		timeStr, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: first argument must be a string")
		}
		
		// Если формат не указан, используем RFC3339
		format := time.RFC3339
		if len(args) == 2 {
			if formatStr, ok := args[1].Any().(string); ok {
				// Конвертируем упрощенный формат в Go формат
				format = convertToGoTimeFormat(formatStr)
			}
		}
		
		parsedTime, err := time.Parse(format, timeStr)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error parsing time: %v", err))
		}
		
		return value.NewTime(parsedTime)
	}
	globalScope.Set("timeFromString", value.NewValue(timeFromStringFunc))

	// timeFormat - форматирование времени
	timeFormatFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeFormat() requires 2 arguments (time, format)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		formatStr, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: second argument must be a format string")
		}
		
		// Конвертируем упрощенный формат в Go формат
		goFormat := convertToGoTimeFormat(formatStr)
		return value.NewString(timeVal.Format(goFormat))
	}
	globalScope.Set("timeFormat", value.NewValue(timeFormatFunc))

	// timeYear - получить год
	timeYearFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: timeYear() requires 1 argument (time)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: argument must be a time value")
		}
		
		return value.NewInt64(int64(timeVal.Year()))
	}
	globalScope.Set("timeYear", value.NewValue(timeYearFunc))

	// timeMonth - получить месяц
	timeMonthFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: timeMonth() requires 1 argument (time)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: argument must be a time value")
		}
		
		return value.NewInt64(int64(timeVal.Month()))
	}
	globalScope.Set("timeMonth", value.NewValue(timeMonthFunc))

	// timeDay - получить день
	timeDayFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: timeDay() requires 1 argument (time)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: argument must be a time value")
		}
		
		return value.NewInt64(int64(timeVal.Day()))
	}
	globalScope.Set("timeDay", value.NewValue(timeDayFunc))

	// timeHour - получить час
	timeHourFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: timeHour() requires 1 argument (time)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: argument must be a time value")
		}
		
		return value.NewInt64(int64(timeVal.Hour()))
	}
	globalScope.Set("timeHour", value.NewValue(timeHourFunc))

	// timeMinute - получить минуты
	timeMinuteFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: timeMinute() requires 1 argument (time)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: argument must be a time value")
		}
		
		return value.NewInt64(int64(timeVal.Minute()))
	}
	globalScope.Set("timeMinute", value.NewValue(timeMinuteFunc))

	// timeSecond - получить секунды
	timeSecondFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: timeSecond() requires 1 argument (time)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: argument must be a time value")
		}
		
		return value.NewInt64(int64(timeVal.Second()))
	}
	globalScope.Set("timeSecond", value.NewValue(timeSecondFunc))

	// timeWeekday - получить день недели (0 = Sunday, 6 = Saturday)
	timeWeekdayFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: timeWeekday() requires 1 argument (time)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: argument must be a time value")
		}
		
		return value.NewInt64(int64(timeVal.Weekday()))
	}
	globalScope.Set("timeWeekday", value.NewValue(timeWeekdayFunc))

	// timeUnix - получить Unix timestamp (секунды)
	timeUnixFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: timeUnix() requires 1 argument (time)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: argument must be a time value")
		}
		
		return value.NewInt64(timeVal.Unix())
	}
	globalScope.Set("timeUnix", value.NewValue(timeUnixFunc))

	// timeAddDays - добавить дни
	timeAddDaysFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeAddDays() requires 2 arguments (time, days)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		days, ok := args[1].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[1].Any().(float64); ok {
				days = int64(floatVal)
			} else {
				return value.NewString("Error: second argument must be a number")
			}
		}
		
		return value.NewTime(timeVal.AddDate(0, 0, int(days)))
	}
	globalScope.Set("timeAddDays", value.NewValue(timeAddDaysFunc))

	// timeAddMonths - добавить месяцы
	timeAddMonthsFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeAddMonths() requires 2 arguments (time, months)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		months, ok := args[1].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[1].Any().(float64); ok {
				months = int64(floatVal)
			} else {
				return value.NewString("Error: second argument must be a number")
			}
		}
		
		return value.NewTime(timeVal.AddDate(0, int(months), 0))
	}
	globalScope.Set("timeAddMonths", value.NewValue(timeAddMonthsFunc))

	// timeAddYears - добавить годы
	timeAddYearsFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeAddYears() requires 2 arguments (time, years)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		years, ok := args[1].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[1].Any().(float64); ok {
				years = int64(floatVal)
			} else {
				return value.NewString("Error: second argument must be a number")
			}
		}
		
		return value.NewTime(timeVal.AddDate(int(years), 0, 0))
	}
	globalScope.Set("timeAddYears", value.NewValue(timeAddYearsFunc))

	// timeAddHours - добавить часы
	timeAddHoursFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeAddHours() requires 2 arguments (time, hours)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		hours, ok := args[1].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[1].Any().(float64); ok {
				hours = int64(floatVal)
			} else {
				return value.NewString("Error: second argument must be a number")
			}
		}
		
		return value.NewTime(timeVal.Add(time.Duration(hours) * time.Hour))
	}
	globalScope.Set("timeAddHours", value.NewValue(timeAddHoursFunc))

	// timeAddMinutes - добавить минуты
	timeAddMinutesFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeAddMinutes() requires 2 arguments (time, minutes)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		minutes, ok := args[1].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[1].Any().(float64); ok {
				minutes = int64(floatVal)
			} else {
				return value.NewString("Error: second argument must be a number")
			}
		}
		
		return value.NewTime(timeVal.Add(time.Duration(minutes) * time.Minute))
	}
	globalScope.Set("timeAddMinutes", value.NewValue(timeAddMinutesFunc))

	// timeAddSeconds - добавить секунды
	timeAddSecondsFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeAddSeconds() requires 2 arguments (time, seconds)")
		}
		
		timeVal, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		seconds, ok := args[1].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[1].Any().(float64); ok {
				seconds = int64(floatVal)
			} else {
				return value.NewString("Error: second argument must be a number")
			}
		}
		
		return value.NewTime(timeVal.Add(time.Duration(seconds) * time.Second))
	}
	globalScope.Set("timeAddSeconds", value.NewValue(timeAddSecondsFunc))

	// timeDiff - разница между двумя временами в секундах
	timeDiffFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeDiff() requires 2 arguments (time1, time2)")
		}
		
		time1, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		time2, ok := args[1].Any().(time.Time)
		if !ok {
			return value.NewString("Error: second argument must be a time value")
		}
		
		diff := time1.Sub(time2)
		return value.NewFloat64(diff.Seconds())
	}
	globalScope.Set("timeDiff", value.NewValue(timeDiffFunc))

	// timeDiffDays - разница в днях
	timeDiffDaysFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeDiffDays() requires 2 arguments (time1, time2)")
		}
		
		time1, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		time2, ok := args[1].Any().(time.Time)
		if !ok {
			return value.NewString("Error: second argument must be a time value")
		}
		
		diff := time1.Sub(time2)
		return value.NewFloat64(diff.Hours() / 24)
	}
	globalScope.Set("timeDiffDays", value.NewValue(timeDiffDaysFunc))

	// timeDiffHours - разница в часах
	timeDiffHoursFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeDiffHours() requires 2 arguments (time1, time2)")
		}
		
		time1, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		time2, ok := args[1].Any().(time.Time)
		if !ok {
			return value.NewString("Error: second argument must be a time value")
		}
		
		diff := time1.Sub(time2)
		return value.NewFloat64(diff.Hours())
	}
	globalScope.Set("timeDiffHours", value.NewValue(timeDiffHoursFunc))

	// timeDiffMinutes - разница в минутах
	timeDiffMinutesFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeDiffMinutes() requires 2 arguments (time1, time2)")
		}
		
		time1, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		time2, ok := args[1].Any().(time.Time)
		if !ok {
			return value.NewString("Error: second argument must be a time value")
		}
		
		diff := time1.Sub(time2)
		return value.NewFloat64(diff.Minutes())
	}
	globalScope.Set("timeDiffMinutes", value.NewValue(timeDiffMinutesFunc))

	// timeBefore - проверка, что time1 < time2
	timeBeforeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeBefore() requires 2 arguments (time1, time2)")
		}
		
		time1, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		time2, ok := args[1].Any().(time.Time)
		if !ok {
			return value.NewString("Error: second argument must be a time value")
		}
		
		return value.NewBool(time1.Before(time2))
	}
	globalScope.Set("timeBefore", value.NewValue(timeBeforeFunc))

	// timeAfter - проверка, что time1 > time2
	timeAfterFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeAfter() requires 2 arguments (time1, time2)")
		}
		
		time1, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		time2, ok := args[1].Any().(time.Time)
		if !ok {
			return value.NewString("Error: second argument must be a time value")
		}
		
		return value.NewBool(time1.After(time2))
	}
	globalScope.Set("timeAfter", value.NewValue(timeAfterFunc))

	// timeEqual - проверка равенства времен
	timeEqualFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: timeEqual() requires 2 arguments (time1, time2)")
		}
		
		time1, ok := args[0].Any().(time.Time)
		if !ok {
			return value.NewString("Error: first argument must be a time value")
		}
		
		time2, ok := args[1].Any().(time.Time)
		if !ok {
			return value.NewString("Error: second argument must be a time value")
		}
		
		return value.NewBool(time1.Equal(time2))
	}
	globalScope.Set("timeEqual", value.NewValue(timeEqualFunc))
}

// convertToGoTimeFormat конвертирует упрощенный формат в Go формат времени
func convertToGoTimeFormat(format string) string {
	// Простые форматы
	switch format {
	case "date":
		return "2006-01-02"
	case "time":
		return "15:04:05"
	case "datetime":
		return "2006-01-02 15:04:05"
	case "rfc3339":
		return time.RFC3339
	case "rfc822":
		return time.RFC822
	case "iso8601":
		return "2006-01-02T15:04:05Z07:00"
	}
	
	// Пользовательский формат - заменяем токены
	result := format
	
	// Год
	result = replaceAll(result, "YYYY", "2006")
	result = replaceAll(result, "YY", "06")
	
	// Месяц
	result = replaceAll(result, "MM", "01")
	result = replaceAll(result, "M", "1")
	result = replaceAll(result, "Mon", "Jan")
	result = replaceAll(result, "Month", "January")
	
	// День
	result = replaceAll(result, "DD", "02")
	result = replaceAll(result, "D", "2")
	
	// Час
	result = replaceAll(result, "HH", "15") // 24-часовой формат
	result = replaceAll(result, "hh", "03") // 12-часовой формат
	result = replaceAll(result, "h", "3")
	
	// Минуты
	result = replaceAll(result, "mm", "04")
	result = replaceAll(result, "m", "4")
	
	// Секунды
	result = replaceAll(result, "ss", "05")
	result = replaceAll(result, "s", "5")
	
	// AM/PM
	result = replaceAll(result, "AM", "PM")
	result = replaceAll(result, "am", "pm")
	
	return result
}

// replaceAll заменяет все вхождения old на new в строке s
func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}