import React, { useState, useEffect } from 'react';
import { Tag, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import CalorieSummary from './CalorieSummary';

interface MealItem {
  name: string;
  quantity: number;
  unit: string;
  calories: number;
}

interface Meal {
  items: MealItem[];
  calories: number;
}

interface MealPlan {
  [key: string]: Meal;
}

interface IngredientCategory {
  name: string;
  ingredients: string[];
}

interface MealIngredients {
  mealName: string;
  categories: IngredientCategory[];
}

const MealPlanner: React.FC = () => {
  const [ingredients, setIngredients] = useState<string[]>([]);
  const [availableIngredients, setAvailableIngredients] = useState<MealIngredients[]>([]);
  const [calories, setCalories] = useState('2000');
  const [mealPlan, setMealPlan] = useState<MealPlan | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchIngredients();
  }, []);

  const fetchIngredients = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/ingredients');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setAvailableIngredients(data);
    } catch (err) {
      console.error('Error fetching ingredients:', err);
      setError('Errore nel caricamento degli ingredienti');
      setAvailableIngredients([]);
    }
  };

  const handleAddIngredient = (value: string) => {
    if (value && !ingredients.includes(value)) {
      setIngredients(prev => [...prev, value]);
    }
  };

  const removeIngredient = (index: number) => {
    setIngredients(ingredients.filter((_, i) => i !== index));
  };

  const generateMealPlan = async () => {
    if (ingredients.length === 0) {
      setError('Seleziona almeno un ingrediente');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const payload = {
        ingredients,
        targetCalories: parseInt(calories)
      };

      const response = await fetch('http://localhost:8080/api/generate-plan', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      setMealPlan(data);
    } catch (err) {
      console.error('Error generating meal plan:', err);
      setError('Errore nella generazione del piano alimentare');
      setMealPlan(null);
    } finally {
      setLoading(false);
    }
  };

  const getIngredientDisplayName = (ingredientKey: string) => {
    for (const meal of availableIngredients) {
      for (const category of meal.categories) {
        const ingredient = category.ingredients.find(i => i === ingredientKey);
        if (ingredient) {
          return ingredient
            .charAt(0).toUpperCase()
            + ingredient.slice(1)
              .replace(/_/g, ' ')
              .trim();
        }
      }
    }
    return ingredientKey;
  };

  return (
    <div className="max-w-6xl mx-auto p-4">
      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Piano Alimentare Personalizzato</CardTitle>
        </CardHeader>
        <CardContent>
          {error && (
            <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
              {error}
            </div>
          )}

          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Ingredienti disponibili
            </label>
            <div className="flex gap-2">
              <div className="flex-1">
                <Select onValueChange={handleAddIngredient}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Seleziona un ingrediente..." />
                  </SelectTrigger>
                  <SelectContent className="max-h-[300px] overflow-y-auto">
                    {availableIngredients.map((meal) => (
                      <SelectGroup key={meal.mealName}>
                        <SelectLabel className="font-bold text-lg py-2 px-2 bg-gray-50">
                          {meal.mealName}
                        </SelectLabel>
                        {meal.categories.map((category) => (
                          <div key={category.name} className="py-1">
                            <SelectLabel className="text-sm font-semibold text-gray-600 px-3">
                              {category.name}
                            </SelectLabel>
                            {category.ingredients.map((ingredient) => (
                              <SelectItem
                                key={ingredient}
                                value={ingredient}
                                className="pl-6"
                                disabled={ingredients.includes(ingredient)}
                              >
                                {getIngredientDisplayName(ingredient)}
                              </SelectItem>
                            ))}
                          </div>
                        ))}
                      </SelectGroup>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <Button
                variant="outline"
                onClick={() => setIngredients([])}
                disabled={ingredients.length === 0}
              >
                Reset
              </Button>
            </div>
          </div>

          {ingredients.length > 0 && (
            <div className="flex flex-wrap gap-2 mb-6 p-4 bg-gray-50 rounded-lg">
              {ingredients.map((ing, index) => (
                <span
                  key={index}
                  className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800"
                >
                  <Tag className="w-4 h-4 mr-1" />
                  {getIngredientDisplayName(ing)}
                  <button
                    onClick={() => removeIngredient(index)}
                    className="ml-2 text-blue-600 hover:text-blue-800"
                  >
                    <X className="w-4 h-4" />
                  </button>
                </span>
              ))}
            </div>
          )}

          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Calorie giornaliere
            </label>
            <input
              type="number"
              value={calories}
              onChange={(e) => setCalories(e.target.value)}
              className="w-40 rounded-md border border-gray-300 px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
              min="1000"
              max="5000"
              step="100"
            />
          </div>

          <Button
            onClick={generateMealPlan}
            disabled={ingredients.length === 0 || loading}
            className="w-full"
          >
            {loading ? (
              <span className="flex items-center justify-center">
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Generazione in corso...
              </span>
            ) : (
              'Genera Piano Alimentare'
            )}
          </Button>
        </CardContent>
      </Card>

      {mealPlan && Object.keys(mealPlan).length > 0 && (
        <>
          <CalorieSummary mealPlan={mealPlan} targetCalories={parseInt(calories)} />

          <div className="mt-6 grid md:grid-cols-2 lg:grid-cols-3 gap-4">
            {Object.entries(mealPlan).map(([mealType, meal]) => (
              <Card key={mealType}>
                <CardHeader>
                  <CardTitle className="text-lg capitalize flex justify-between items-center">
                    <span>{mealType}</span>
                    <span className="text-sm font-normal text-gray-600">
                      {meal.calories} kcal
                    </span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    {meal.items && meal.items.map((item: MealItem, idx: number) => (
                      <div key={idx} className="flex justify-between text-sm border-b border-gray-200 py-2 last:border-0">
                        <span className="font-medium">{item.name}</span>
                        <div className="text-gray-600">
                          <span>{item.quantity} {item.unit}</span>
                          <span className="ml-2">({item.calories} kcal)</span>
                        </div>
                      </div>
                    ))}
                    {(!meal.items || meal.items.length === 0) && (
                      <div className="text-sm text-gray-500 italic">
                        Nessun alimento selezionato per questo pasto
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </>
      )}
    </div>
  );
};

export default MealPlanner;