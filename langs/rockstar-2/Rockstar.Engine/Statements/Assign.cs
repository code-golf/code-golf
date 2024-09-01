using System.Text;
using Rockstar.Engine.Expressions;
using Rockstar.Engine.Values;

namespace Rockstar.Engine.Statements;

public class Declare(Variable variable, Expression? expression = null) : Statement {
	public Variable Variable => variable;
	public Expression? Expression => expression;

	public override StringBuilder Print(StringBuilder sb, string prefix) {
		sb.AppendLine(prefix + "declare:");
		variable.Print(sb, prefix + INDENT);
		return (expression == default ? sb : expression.Print(sb, prefix + INDENT));
	}
}

public class Assign(Variable variable, Expression expression) : ExpressionStatement(expression) {
	public Variable Variable => variable;
	public override StringBuilder Print(StringBuilder sb, string prefix) {
		sb.AppendLine(prefix + "assign:");
		variable.Print(sb, prefix + INDENT);
		return Expression.Print(sb, prefix + INDENT);
	}
}

public class Ninja(Variable variable, Number number) : Statement {
	public Variable Variable => variable;
	public Number Number => number;
	public override StringBuilder Print(StringBuilder sb, string prefix) {
		sb.AppendLine(prefix + "ninja:");
		variable.Print(sb, prefix + INDENT);
		return number.Print(sb, prefix + INDENT);
	}
}

public class Listen : Statement {
	public Variable? Variable { get; init; } = default;
	public Listen() {}
	public Listen(Variable variable) => this.Variable = variable;
	public override StringBuilder Print(StringBuilder sb, string prefix) {
		sb.Append(prefix).Append("listen:");
		if (Variable != default) sb.Append(" => ").Append(Variable.Name);
		return sb.AppendLine();
	}
}
