/* ************************************************************************** */
/*                                                                            */
/*                                                        ::::::::            */
/*   main.c                                             :+:    :+:            */
/*                                                     +:+                    */
/*   By: jkoers <jkoers@student.codam.nl>             +#+                     */
/*                                                   +#+                      */
/*   Created: 2020/10/26 17:30:43 by jkoers        #+#    #+#                 */
/*   Updated: 2022/11/18 09:35:58 by jkoers        ########   odam.nl         */
/*                                                                            */
/* ************************************************************************** */

#include <stddef.h>
#include <stdbool.h>
#include <stdio.h>

static int ft_isdigit(int c)
{
	return (c >= '0' && c <= '9');
}


static bool	ft_isspace(char c)
{
	if (c == ' ')
		return (true);
	if (c == '\t')
		return (true);
	if (c == '\v')
		return (true);
	if (c == '\f')
		return (true);
	if (c == '\r')
		return (true);
	return (false);
}

int			ft_atoi(char *str)
{
	int		result;
	bool	is_negative;

	while (ft_isspace(*str))
		str++;
	is_negative = *str == '-';
	if (*str == '-' || *str == '+')
		str++;
	result = 0;
	while (ft_isdigit(*str))
	{
		result *= 10;
		result -= (int)(*str - 47);
		str++;
	}
	return (is_negative ? result : (-result));
}

int main() {
	printf("%d\n", ft_atoi("123"));
	printf("%d\n", ft_atoi("42"));
	printf("%d\n", ft_atoi("0"));
	printf("%d\n", ft_atoi("-42"));
	printf("%d\n", ft_atoi("2147483647"));
	printf("%d\n", ft_atoi("-2147483648"));
	printf("%d\n", ft_atoi("  +1"));
	printf("%d\n", ft_atoi("  -42"));
}
