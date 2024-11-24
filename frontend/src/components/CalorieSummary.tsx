// src/components/CalorieSummary.tsx
import React from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';

const CalorieSummary = ({ mealPlan, targetCalories }) => {
    if (!mealPlan) return null;

    const meals = Object.entries(mealPlan).map(([name, meal]) => ({
        name: name.charAt(0).toUpperCase() + name.slice(1),
        calories: meal.calories,
        expectedPercentage:
            name === 'colazione' ? 25 :
                name === 'spuntino' ? 10 :
                    name === 'pranzo' ? 35 :
                        name === 'merenda' ? 10 :
                            name === 'cena' ? 20 : 0,
        expectedCalories: (targetCalories * (
            name === 'colazione' ? 0.25 :
                name === 'spuntino' ? 0.10 :
                    name === 'pranzo' ? 0.35 :
                        name === 'merenda' ? 0.10 :
                            name === 'cena' ? 0.20 : 0
        ))
    }));

    const totalCalories = meals.reduce((sum, meal) => sum + meal.calories, 0);
    const calorieDeficit = targetCalories - totalCalories;

    // Colors for the pie chart
    const COLORS = ['#FF8042', '#00C49F', '#0088FE', '#FFBB28', '#FF6B6B'];

    return (
        <Card className="mt-6">
            <CardHeader>
                <CardTitle className="flex justify-between items-center">
                    <span>Riepilogo Calorico</span>
                    <span className="text-lg font-normal">
                        {totalCalories} / {targetCalories} kcal
                    </span>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="grid md:grid-cols-2 gap-6">
                    <div className="h-64">
                        <ResponsiveContainer width="100%" height="100%">
                            <PieChart>
                                <Pie
                                    data={meals}
                                    dataKey="calories"
                                    nameKey="name"
                                    cx="50%"
                                    cy="50%"
                                    outerRadius={80}
                                    fill="#8884d8"
                                    label={({ name, calories }) => `${name}: ${calories}kcal`}
                                >
                                    {meals.map((entry, index) => (
                                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                    ))}
                                </Pie>
                                <Tooltip />
                            </PieChart>
                        </ResponsiveContainer>
                    </div>

                    <div className="space-y-4">
                        {meals.map((meal, index) => (
                            <div key={meal.name} className="flex justify-between items-center">
                                <div className="flex items-center">
                                    <div
                                        className="w-3 h-3 rounded-full mr-2"
                                        style={{ backgroundColor: COLORS[index % COLORS.length] }}
                                    />
                                    <span className="font-medium">{meal.name}</span>
                                </div>
                                <div className="text-right">
                                    <div>{meal.calories} / {Math.round(meal.expectedCalories)} kcal</div>
                                    <div className="text-sm text-gray-500">
                                        {((meal.calories / targetCalories) * 100).toFixed(1)}% / {meal.expectedPercentage}%
                                    </div>
                                </div>
                            </div>
                        ))}

                        {Math.abs(calorieDeficit) > 50 && (
                            <div className={`mt-4 p-3 rounded-lg text-sm ${calorieDeficit > 0 ? 'bg-yellow-50 text-yellow-800' : 'bg-red-50 text-red-800'}`}>
                                {calorieDeficit > 0
                                    ? `Mancano ${Math.round(calorieDeficit)} kcal per raggiungere l'obiettivo giornaliero`
                                    : `Il piano supera di ${Math.abs(Math.round(calorieDeficit))} kcal l'obiettivo giornaliero`
                                }
                            </div>
                        )}
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};

export default CalorieSummary;